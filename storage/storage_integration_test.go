package storage

import (
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/testutil"
	"github.com/KonnorFrik/getman/types"
)

func TestIntegrationFileStorage_FullCycle(t *testing.T) {
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

	collectionsDir := fs.CollectionsDir()
	if collectionsDir == "" {
		t.Fatal("expected collections directory to be set")
	}

	environmentsDir := fs.EnvironmentsDir()
	if environmentsDir == "" {
		t.Fatal("expected environments directory to be set")
	}

	historyDir := fs.HistoryDir()
	if historyDir == "" {
		t.Fatal("expected history directory to be set")
	}

	logsDir := fs.LogsDir()
	if logsDir == "" {
		t.Fatal("expected logs directory to be set")
	}

	configPath := fs.ConfigPath()
	if configPath == "" {
		t.Fatal("expected config path to be set")
	}
}

func TestIntegrationHistoryStorage_FullCycle(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hs := NewHistoryStorage(fs)

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
				Response: &types.Response{
					StatusCode: 200,
					Status:     "200 OK",
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

	list, err := hs.List()
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}

	if len(list) == 0 {
		t.Fatal("expected history to contain items")
	}

	last, err := hs.GetLast()
	if err != nil {
		t.Fatalf("unexpected error getting last: %v", err)
	}

	if last.CollectionName != result.CollectionName {
		t.Errorf("expected collection name %s, got %s", result.CollectionName, last.CollectionName)
	}

	history, err := hs.GetHistory(10)
	if err != nil {
		t.Fatalf("unexpected error getting history: %v", err)
	}

	if len(history) == 0 {
		t.Fatal("expected history to contain items")
	}

	if err := hs.Clear(); err != nil {
		t.Fatalf("unexpected error clearing: %v", err)
	}

	list, err = hs.List()
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("expected history to be empty, got %d items", len(list))
	}
}

func TestIntegrationLogStorage_FullCycle(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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
		t.Fatalf("unexpected error saving: %v", err)
	}

	last, err := ls.GetLast()
	if err != nil {
		t.Fatalf("unexpected error getting last: %v", err)
	}

	if len(last) == 0 {
		t.Fatal("expected logs to contain data")
	}
}

func TestIntegrationStorage_ConcurrentAccess(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	fs, err := NewFileStorage(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hs := NewHistoryStorage(fs)

	var wg sync.WaitGroup
	numGoroutines := 2
	numOperations := 5

	wg.Add(numGoroutines)

	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for range numOperations {
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

				time.Sleep(time.Duration(rand.IntN(100)) * time.Millisecond)
				hs.Save(result)
			}
		}(i)
	}

	wg.Wait()

	list, err := hs.List()
	if err != nil {
		t.Fatalf("unexpected error listing: %v", err)
	}

	if len(list) != numGoroutines*numOperations {
		t.Errorf("expected %d items, got %d", numGoroutines*numOperations, len(list))
	}
}

