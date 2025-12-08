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

	"github.com/KonnorFrik/getman/errors"
	"github.com/KonnorFrik/getman/types"
)

type HistoryStorage struct {
	fileStorage *FileStorage
}

func NewHistoryStorage(fileStorage *FileStorage) *HistoryStorage {
	return &HistoryStorage{
		fileStorage: fileStorage,
	}
}

// Save - save 'result' into file with path "<base_path>/<getman_dir>/history/<timestamp>.json
// If file exist - it will be overwrited.
func (hs *HistoryStorage) Save(result *types.ExecutionResult) error {
	timestamp := FormatTimestamp(time.Now())
	filename := fmt.Sprintf("%s.json", timestamp)

	dir := hs.fileStorage.HistoryDir()
	filePath := filepath.Join(dir, filename)
	data, err := json.MarshalIndent(result, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal execution result: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("%w: failed to write history file: %w", errors.ErrStorageError, err)
	}

	return nil
}

func (hs *HistoryStorage) Load(timestamp string) (*types.ExecutionResult, error) {
	filename := fmt.Sprintf("%s.json", timestamp)
	filePath := filepath.Join(hs.fileStorage.HistoryDir(), filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	var result types.ExecutionResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse history file: %w", err)
	}

	return &result, nil
}

func (hs *HistoryStorage) List() ([]string, error) {
	dir := hs.fileStorage.HistoryDir()
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil, fmt.Errorf("failed to read history directory: %w", err)
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

	sort.Sort(sort.Reverse(sort.StringSlice(timestamps)))
	return timestamps, nil
}

func (hs *HistoryStorage) GetLast() (*types.ExecutionResult, error) {
	timestamps, err := hs.List()
	if err != nil {
		return nil, err
	}

	if len(timestamps) == 0 {
		return nil, fmt.Errorf("no history files found")
	}

	return hs.Load(timestamps[0])
}

func (hs *HistoryStorage) GetHistory(limit int) ([]*types.RequestExecution, error) {
	timestamps, err := hs.List()

	if err != nil {
		return nil, err
	}

	if limit > len(timestamps) {
		limit = len(timestamps)
	}

	var allExecutions []*types.RequestExecution

	for i := 0; i < limit; i++ {
		result, err := hs.Load(timestamps[i])

		if err != nil {
			continue
		}

		allExecutions = append(allExecutions, result.Requests...)
	}

	return allExecutions, nil
}

func (hs *HistoryStorage) Clear() error {
	dir := hs.fileStorage.HistoryDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read history directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to remove history file %s: %w", filePath, err)
		}
	}

	return nil
}
