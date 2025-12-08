/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package helper

import (
	"os"
	"path/filepath"
)

func CreateTempDir() (string, error) {
	return os.MkdirTemp("", "getman_test_*")
}

func CleanupTempDir(dir string) error {
	return os.RemoveAll(dir)
}

func WriteTestFile(dir, filename, content string) error {
	filePath := filepath.Join(dir, filename)
	return os.WriteFile(filePath, []byte(content), 0644)
}

func ReadTestFile(dir, filename string) (string, error) {
	filePath := filepath.Join(dir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
