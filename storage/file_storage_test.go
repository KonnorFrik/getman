package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/testutil"
)

func TestUnitNewFileStorage_ValidPath(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fs == nil {
		t.Fatal("expected file storage to be created")
	}

	if fs.basePath != dir {
		t.Errorf("expected base path %s, got %s", dir, fs.basePath)
	}
}

func TestUnitNewFileStorage_HomeDirPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	fs, err := NewFileStorage("~/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := filepath.Join(homeDir, "test")
	if fs.basePath != expectedPath {
		t.Errorf("expected base path %s, got %s", expectedPath, fs.basePath)
	}

	testutil.CleanupTempDir(expectedPath)
}

func TestUnitNewFileStorage_CreateDirectories(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collectionsDir := fs.CollectionsDir()
	if _, err := os.Stat(collectionsDir); os.IsNotExist(err) {
		t.Errorf("expected collections directory to be created: %s", collectionsDir)
	}

	environmentsDir := fs.EnvironmentsDir()
	if _, err := os.Stat(environmentsDir); os.IsNotExist(err) {
		t.Errorf("expected environments directory to be created: %s", environmentsDir)
	}

	historyDir := fs.HistoryDir()
	if _, err := os.Stat(historyDir); os.IsNotExist(err) {
		t.Errorf("expected history directory to be created: %s", historyDir)
	}

	logsDir := fs.LogsDir()
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		t.Errorf("expected logs directory to be created: %s", logsDir)
	}
}

func TestUnitBasePath(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fs.BasePath() != dir {
		t.Errorf("expected base path %s, got %s", dir, fs.BasePath())
	}
}

func TestUnitCollectionsDir(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDir := filepath.Join(dir, "collections")
	if fs.CollectionsDir() != expectedDir {
		t.Errorf("expected collections dir %s, got %s", expectedDir, fs.CollectionsDir())
	}
}

func TestUnitEnvironmentsDir(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDir := filepath.Join(dir, "environments")
	if fs.EnvironmentsDir() != expectedDir {
		t.Errorf("expected environments dir %s, got %s", expectedDir, fs.EnvironmentsDir())
	}
}

func TestUnitHistoryDir(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDir := filepath.Join(dir, "history")
	if fs.HistoryDir() != expectedDir {
		t.Errorf("expected history dir %s, got %s", expectedDir, fs.HistoryDir())
	}
}

func TestUnitLogsDir(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedDir := filepath.Join(dir, "logs")
	if fs.LogsDir() != expectedDir {
		t.Errorf("expected logs dir %s, got %s", expectedDir, fs.LogsDir())
	}
}

func TestUnitConfigPath(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPath := filepath.Join(dir, "config.yaml")
	if fs.ConfigPath() != expectedPath {
		t.Errorf("expected config path %s, got %s", expectedPath, fs.ConfigPath())
	}
}

func TestUnitFormatTimestamp(t *testing.T) {
	timestamp := time.Date(2025, 12, 1, 22, 55, 39, 0, time.UTC)
	formatted := FormatTimestamp(timestamp)
	expected := "01_12_25_22_55_39.0000"

	if formatted != expected {
		t.Errorf("expected formatted timestamp %s, got %s", expected, formatted)
	}
}

func TestUnitParseTimestamp(t *testing.T) {
	timestampStr := "01_12_25_22_55_39.0000"
	timestamp, err := ParseTimestamp(timestampStr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2025, 12, 1, 22, 55, 39, 0, time.UTC)
	if !timestamp.Equal(expected) {
		t.Errorf("expected timestamp %v, got %v", expected, timestamp)
	}
}

func TestUnitParseTimestamp_Invalid(t *testing.T) {
	_, err := ParseTimestamp("invalid")
	if err == nil {
		t.Fatal("expected error for invalid timestamp")
	}
}

func TestUnitExpandPath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home dir: %v", err)
	}

	expanded, err := ExpandPath("~/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(homeDir, "test")
	if expanded != expected {
		t.Errorf("expected expanded path %s, got %s", expected, expanded)
	}
}

func TestUnitExpandPath_NoTilde(t *testing.T) {
	path := "/absolute/path"
	expanded, err := ExpandPath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if expanded != path {
		t.Errorf("expected path %s, got %s", path, expanded)
	}
}

// func TestUnitFormatTimestamp_RoundTrip(t *testing.T) {
// 	original := time.Now()
// 	formatted := FormatTimestamp(original)
// 	parsed, err := ParseTimestamp(formatted)
//
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
//
// 	origTrunc := original.In(time.FixedZone("+05:00", 5*60*60)).Truncate(time.Second)
//
// 	if !parsed.Equal(origTrunc) {
// 		t.Errorf("expected parsed timestamp to match original (within second), got %v vs %v", parsed, origTrunc)
// 	}
// }
