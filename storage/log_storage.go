/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/KonnorFrik/getman/types"
)

// LogStorage provides storage for log entries.
type LogStorage struct {
	fileStorage *FileStorage
}

// NewLogStorage creates a new LogStorage instance.
func NewLogStorage(fileStorage *FileStorage) *LogStorage {
	return &LogStorage{
		fileStorage: fileStorage,
	}
}

// Save saves log entries to a timestamped JSON file.
func (ls *LogStorage) Save(logs []types.LogEntry) error {
	timestamp := FormatTimestamp(time.Now())
	filename := fmt.Sprintf("%s.json", timestamp)
	filePath := filepath.Join(ls.fileStorage.LogsDir(), filename)

	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal logs: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write log file: %w", err)
	}

	return nil
}

// Load loads log entries by timestamp.
func (ls *LogStorage) Load(timestamp string) ([]types.LogEntry, error) {
	filename := fmt.Sprintf("%s.json", timestamp)
	filePath := filepath.Join(ls.fileStorage.LogsDir(), filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	var logs []types.LogEntry
	if err := json.Unmarshal(data, &logs); err != nil {
		return nil, fmt.Errorf("failed to parse log file: %w", err)
	}

	return logs, nil
}

// GetLast returns the most recent log entries as JSON bytes.
func (ls *LogStorage) GetLast() ([]byte, error) {
	dir := ls.fileStorage.LogsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs directory: %w", err)
	}

	var timestamps []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if len(name) > 5 && name[len(name)-5:] == ".json" {
			timestamp := name[:len(name)-5]
			if _, err := ParseTimestamp(timestamp); err == nil {
				timestamps = append(timestamps, timestamp)
			}
		}
	}

	if len(timestamps) == 0 {
		return nil, fmt.Errorf("no log files found")
	}

	sort.Sort(sort.Reverse(sort.StringSlice(timestamps)))

	logs, err := ls.Load(timestamps[0])
	if err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal logs: %w", err)
	}

	return data, nil
}
