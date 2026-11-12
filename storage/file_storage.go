/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FileStorage provides file system-based storage for collections, environments, history, and logs.
type FileStorage struct {
	basePath string
}

// NewFileStorage creates a new FileStorage instance with the specified base path.
func NewFileStorage(basePath string) (*FileStorage, error) {
	expandedPath, err := expandPath(basePath)

	if err != nil {
		return nil, fmt.Errorf("failed to expand path: %w", err)
	}

	fs := &FileStorage{
		basePath: expandedPath,
	}

	if err := fs.ensureDirectories(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	return fs, nil
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()

		if err != nil {
			return "", err
		}

		result := filepath.Join(homeDir, path[2:])
		return result, nil
	}

	return path, nil
}

// ExpandPath expands a path string, replacing ~ with the user's home directory.
func ExpandPath(path string) (string, error) {
	return expandPath(path)
}

func (fs *FileStorage) ensureDirectories() error {
	dirs := []string{
		fs.CollectionsDir(),
		fs.EnvironmentsDir(),
		fs.HistoryDir(),
		fs.LogsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// BasePath returns the base path of the file storage.
func (fs *FileStorage) BasePath() string {
	return fs.basePath
}

// CollectionsDir returns the path to the collections directory.
func (fs *FileStorage) CollectionsDir() string {
	const dirName = "collections"
	return filepath.Join(fs.basePath, dirName)
}

// EnvironmentsDir returns the path to the environments directory.
func (fs *FileStorage) EnvironmentsDir() string {
	const dirName = "environments"
	return filepath.Join(fs.basePath, dirName)
}

// HistoryDir returns the path to the history directory.
func (fs *FileStorage) HistoryDir() string {
	const dirName = "history"
	return filepath.Join(fs.basePath, dirName)
}

// LogsDir returns the path to the logs directory.
func (fs *FileStorage) LogsDir() string {
	const dirName = "logs"
	return filepath.Join(fs.basePath, dirName)
}

// ConfigPath returns the path to the configuration file.
func (fs *FileStorage) ConfigPath() string {
	const fileName = "config.yaml"
	return filepath.Join(fs.basePath, fileName)
}

var timestampedWriteMu sync.Mutex

// writeTimestampedFile writes data to a new file in dir using the current time as
// the filename. If a file with the same timestamp already exists, the next
// millisecond is used.
func writeTimestampedFile(dir string, data []byte) error {
	timestampedWriteMu.Lock()
	defer timestampedWriteMu.Unlock()

	now := time.Now()

	for range 1000 {
		filePath := filepath.Join(dir, FormatTimestamp(now)+".json")

		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err == nil {
			if _, err := file.Write(data); err != nil {
				file.Close()
				return err
			}

			return file.Close()
		}

		if !os.IsExist(err) {
			return err
		}

		now = now.Add(time.Millisecond)
	}

	return fmt.Errorf("failed to create unique timestamped file in %s", dir)
}

const timeFormat string = "02_01_06_15_04_05.0000"

// FormatTimestamp formats a time value as a string for use in filenames.
func FormatTimestamp(t time.Time) string {
	return t.Format(timeFormat)
}

// ParseTimestamp parses a timestamp string from a filename.
func ParseTimestamp(s string) (time.Time, error) {
	return time.Parse(timeFormat, s)
}
