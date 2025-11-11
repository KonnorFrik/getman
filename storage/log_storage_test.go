package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/testutil"
	"github.com/KonnorFrik/getman/types"
)

func TestUnitNewLogStorage(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)
	if ls == nil {
		t.Fatal("expected log storage to be created")
	}

	if ls.fileStorage != fs {
		t.Error("expected file storage to be set")
	}
}

func TestUnitLogStorage_Save_Valid(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)

	logs := []types.LogEntry{
		{
			Time:    time.Now(),
			Level:   "INFO",
			Message: "Test log message",
		},
		{
			Time:    time.Now(),
			Level:   "ERROR",
			Message: "Test error message",
		},
	}

	if err := ls.Save(logs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := os.ReadDir(fs.LogsDir())
	if err != nil {
		t.Fatalf("unexpected error reading logs dir: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}

func TestUnitLogStorage_Load_Valid(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)

	logs := []types.LogEntry{
		{
			Time:    time.Now(),
			Level:   "INFO",
			Message: "Test log message",
		},
	}

	if err := ls.Save(logs); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	files, err := os.ReadDir(fs.LogsDir())
	if err != nil {
		t.Fatalf("unexpected error reading logs dir: %v", err)
	}

	if len(files) == 0 {
		t.Fatal("expected at least one file")
	}

	filename := files[0].Name()
	timestamp := filename[:len(filename)-5]

	loaded, err := ls.Load(timestamp)
	if err != nil {
		t.Fatalf("unexpected error loading: %v", err)
	}

	if len(loaded) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(loaded))
	}

	if loaded[0].Level != "INFO" {
		t.Errorf("expected level INFO, got %s", loaded[0].Level)
	}

	if loaded[0].Message != "Test log message" {
		t.Errorf("expected message 'Test log message', got %s", loaded[0].Message)
	}
}

func TestUnitLogStorage_Load_NotFound(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)

	_, err = ls.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent timestamp")
	}
}

func TestUnitLogStorage_GetLast_Exists(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)

	logs := []types.LogEntry{
		{
			Time:    time.Now(),
			Level:   "INFO",
			Message: "Test log message",
		},
	}

	if err := ls.Save(logs); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	last, err := ls.GetLast()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(last) == 0 {
		t.Fatal("expected last logs to contain data")
	}

	var loadedLogs []types.LogEntry
	if err := json.Unmarshal(last, &loadedLogs); err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	if len(loadedLogs) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(loadedLogs))
	}
}

func TestUnitLogStorage_GetLast_NotFound(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)

	_, err = ls.GetLast()
	if err == nil {
		t.Fatal("expected error for empty logs")
	}
}

func TestUnitLogStorage_GetLast_Empty(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)

	_, err = ls.GetLast()
	if err == nil {
		t.Fatal("expected error for empty logs")
	}
}

func TestUnitLogStorage_Save_MultipleLogs(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)

	logs := []types.LogEntry{
		{
			Time:    time.Now(),
			Level:   "INFO",
			Message: "Log 1",
		},
		{
			Time:    time.Now(),
			Level:   "WARN",
			Message: "Log 2",
		},
		{
			Time:    time.Now(),
			Level:   "ERROR",
			Message: "Log 3",
		},
	}

	if err := ls.Save(logs); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := os.ReadDir(fs.LogsDir())
	if err != nil {
		t.Fatalf("unexpected error reading logs dir: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}

	filename := files[0].Name()
	timestamp := filename[:len(filename)-5]

	loaded, err := ls.Load(timestamp)
	if err != nil {
		t.Fatalf("unexpected error loading: %v", err)
	}

	if len(loaded) != 3 {
		t.Errorf("expected 3 log entries, got %d", len(loaded))
	}
}

func TestUnitLogStorage_GetLast_MultipleFiles(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)

	for i := 0; i < 3; i++ {
		logs := []types.LogEntry{
			{
				Time:    time.Now(),
				Level:   "INFO",
				Message: fmt.Sprintf("Log %d", i),
			},
		}

		if err := ls.Save(logs); err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	last, err := ls.GetLast()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var loadedLogs []types.LogEntry
	if err := json.Unmarshal(last, &loadedLogs); err != nil {
		t.Fatalf("unexpected error unmarshaling: %v", err)
	}

	if len(loadedLogs) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(loadedLogs))
	}

	if loadedLogs[0].Message != "Log 2" {
		t.Errorf("expected message 'Log 2', got %s", loadedLogs[0].Message)
	}
}

func TestUnitLogStorage_Load_IgnoreInvalidFiles(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	ls := NewLogStorage(fs)

	invalidFile := filepath.Join(fs.LogsDir(), "invalid.txt")
	if err := os.WriteFile(invalidFile, []byte("invalid"), 0644); err != nil {
		t.Fatalf("failed to write invalid file: %v", err)
	}

	_, err = ls.GetLast()
	if err == nil {
		t.Fatal("expected error when only invalid files exist")
	}
}
