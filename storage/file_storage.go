package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileStorage struct {
	basePath string
}

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
		return filepath.Join(homeDir, path[2:]), nil
	}
	return path, nil
}

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

func (fs *FileStorage) BasePath() string {
	return fs.basePath
}

func (fs *FileStorage) CollectionsDir() string {
	return filepath.Join(fs.basePath, "collections")
}

func (fs *FileStorage) EnvironmentsDir() string {
	return filepath.Join(fs.basePath, "environments")
}

func (fs *FileStorage) HistoryDir() string {
	return filepath.Join(fs.basePath, "history")
}

func (fs *FileStorage) LogsDir() string {
	return filepath.Join(fs.basePath, "logs")
}

func (fs *FileStorage) ConfigPath() string {
	return filepath.Join(fs.basePath, "config.yaml")
}

func FormatTimestamp(t time.Time) string {
	return t.Format("02_01_06_15_04_05")
}

func ParseTimestamp(s string) (time.Time, error) {
	return time.Parse("02_01_06_15_04_05", s)
}
