package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/testutil"
	"github.com/KonnorFrik/getman/types"
)

func TestUnitNewHistoryStorage(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)
	if hs == nil {
		t.Fatal("expected history storage to be created")
	}

	if hs.fileStorage != fs {
		t.Error("expected file storage to be set")
	}
}

func TestUnitHistoryStorage_Save_Valid(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	result := &types.ExecutionResult{
		CollectionName: "Test Collection",
		Environment:    "test",
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		TotalDuration:  time.Second,
		Requests:       []*types.RequestExecution{},
		Statistics: &types.Statistics{
			Total:   0,
			Success: 0,
			Failed:  0,
			AvgTime: 0,
			MinTime: 0,
			MaxTime: 0,
		},
	}

	if err := hs.Save(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	files, err := os.ReadDir(fs.HistoryDir())
	if err != nil {
		t.Fatalf("unexpected error reading history dir: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}

func TestUnitHistoryStorage_Load_Valid(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	result := &types.ExecutionResult{
		CollectionName: "Test Collection",
		Environment:    "test",
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		TotalDuration:  time.Second,
		Requests:       []*types.RequestExecution{},
		Statistics: &types.Statistics{
			Total:   0,
			Success: 0,
			Failed:  0,
			AvgTime: 0,
			MinTime: 0,
			MaxTime: 0,
		},
	}

	if err := hs.Save(result); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	timestamps, err := hs.List()
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}

	if len(timestamps) == 0 {
		t.Fatal("expected at least one timestamp")
	}

	loaded, err := hs.Load(timestamps[0])
	if err != nil {
		t.Fatalf("unexpected error loading: %v", err)
	}

	if loaded.CollectionName != result.CollectionName {
		t.Errorf("expected collection name %s, got %s", result.CollectionName, loaded.CollectionName)
	}
}

func TestUnitHistoryStorage_Load_NotFound(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	_, err = hs.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent timestamp")
	}
}

func TestUnitList_Empty(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	list, err := hs.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("expected empty list, got %d items", len(list))
	}
}

func TestUnitList_MultipleFiles(t *testing.T) {
	const historyEniriesLen int = 2
	dir, err := testutil.CreateTempDir()

	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	defer testutil.CleanupTempDir(dir)
	fs, err := NewFileStorage(dir)

	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	for i := range historyEniriesLen {
		result := &types.ExecutionResult{
			CollectionName: "Test Collection",
			Environment:    "test",
			StartTime:      time.Now().Add(time.Duration(i) * time.Second),
			EndTime:        time.Now().Add(time.Duration(i) * time.Second),
			TotalDuration:  time.Second,
			Requests:       []*types.RequestExecution{},
			Statistics: &types.Statistics{
				Total:   0,
				Success: 0,
				Failed:  0,
				AvgTime: 0,
				MinTime: 0,
				MaxTime: 0,
			},
		}

		if err := hs.Save(result); err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	list, err := hs.List()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != historyEniriesLen {
		t.Errorf("expected %d items, got %d", historyEniriesLen, len(list))
	}

	for i := 0; i < len(list)-1; i++ {
		if list[i] <= list[i+1] {
			t.Error("expected list to be sorted in descending order")
		}
	}
}

func TestUnitHistoryStorage_GetLast_Exists(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	result := &types.ExecutionResult{
		CollectionName: "Test Collection",
		Environment:    "test",
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		TotalDuration:  time.Second,
		Requests:       []*types.RequestExecution{},
		Statistics: &types.Statistics{
			Total:   0,
			Success: 0,
			Failed:  0,
			AvgTime: 0,
			MinTime: 0,
			MaxTime: 0,
		},
	}

	if err := hs.Save(result); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	last, err := hs.GetLast()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if last.CollectionName != result.CollectionName {
		t.Errorf("expected collection name %s, got %s", result.CollectionName, last.CollectionName)
	}
}

func TestUnitHistoryStorage_GetLast_NotFound(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	_, err = hs.GetLast()
	if err == nil {
		t.Fatal("expected error for empty history")
	}
}

func TestUnitGetHistory_WithLimit(t *testing.T) {
	const historyEniriesLen int = 2
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	for range historyEniriesLen {
		result := &types.ExecutionResult{
			CollectionName: "Test Collection",
			Environment:    "test",
			StartTime:      time.Now(),
			EndTime:        time.Now(),
			TotalDuration:  time.Second,
			Requests: []*types.RequestExecution{
				{
					Request: &types.Request{
						Method: "GET",
						URL:    "http://example.com",
					},
					Duration:  time.Millisecond * 100,
					Timestamp: time.Now(),
				},
			},
			Statistics: &types.Statistics{
				Total:   1,
				Success: 1,
				Failed:  0,
				AvgTime: time.Millisecond * 100,
				MinTime: time.Millisecond * 100,
				MaxTime: time.Millisecond * 100,
			},
		}

		if err := hs.Save(result); err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	const wantHistoryEntriesLen int = historyEniriesLen - 1
	history, err := hs.GetHistory(wantHistoryEntriesLen)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(history) != wantHistoryEntriesLen {
		t.Errorf("expected %d executions, got %d", wantHistoryEntriesLen, len(history))
	}
}

func TestUnitGetHistory_WithoutLimit(t *testing.T) {
	const historyEniriesLen int = 2
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	for range historyEniriesLen {
		result := &types.ExecutionResult{
			CollectionName: "Test Collection",
			Environment:    "test",
			StartTime:      time.Now(),
			EndTime:        time.Now(),
			TotalDuration:  time.Second,
			Requests: []*types.RequestExecution{
				{
					Request: &types.Request{
						Method: "GET",
						URL:    "http://example.com",
					},
					Duration:  time.Millisecond * 100,
					Timestamp: time.Now(),
				},
			},
			Statistics: &types.Statistics{
				Total:   1,
				Success: 1,
				Failed:  0,
				AvgTime: time.Millisecond * 100,
				MinTime: time.Millisecond * 100,
				MaxTime: time.Millisecond * 100,
			},
		}

		if err := hs.Save(result); err != nil {
			t.Fatalf("unexpected error saving: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
	}

	const wantHistoryEntriesLen int = historyEniriesLen * 2
	history, err := hs.GetHistory(wantHistoryEntriesLen)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(history) != historyEniriesLen {
		t.Errorf("expected %d executions, got %d", historyEniriesLen, len(history))
	}
}

func TestUnitClear(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	result := &types.ExecutionResult{
		CollectionName: "Test Collection",
		Environment:    "test",
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		TotalDuration:  time.Second,
		Requests:       []*types.RequestExecution{},
		Statistics: &types.Statistics{
			Total:   0,
			Success: 0,
			Failed:  0,
			AvgTime: 0,
			MinTime: 0,
			MaxTime: 0,
		},
	}

	if err := hs.Save(result); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	if err := hs.Clear(); err != nil {
		t.Fatalf("unexpected error clearing: %v", err)
	}

	list, err := hs.List()
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("expected empty list after clear, got %d items", len(list))
	}
}

func TestUnitList_IgnoreInvalidFiles(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	hs := NewHistoryStorage(fs)

	invalidFile := filepath.Join(fs.HistoryDir(), "invalid.txt")
	if err := os.WriteFile(invalidFile, []byte("invalid"), 0644); err != nil {
		t.Fatalf("failed to write invalid file: %v", err)
	}

	list, err := hs.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("expected empty list (invalid files should be ignored), got %d items", len(list))
	}
}
