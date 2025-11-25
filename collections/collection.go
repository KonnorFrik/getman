package collections

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/KonnorFrik/getman/storage"
	"github.com/KonnorFrik/getman/types"
)

type Collection struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Items       []*types.RequestItem `json:"items"`
	EnvName    string          `json:"environment_name"`
}

func LoadCollectionFromFile(filePath string) (*Collection, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection file: %w", err)
	}

	var collection Collection
	if err := json.Unmarshal(data, &collection); err != nil {
		return nil, fmt.Errorf("failed to parse collection file: %w", err)
	}

	if err := validateCollection(&collection); err != nil {
		return nil, fmt.Errorf("invalid collection: %w", err)
	}

	return &collection, nil
}

func SaveCollectionToFile(collection *Collection, filePath string) error {
	if err := validateCollection(collection); err != nil {
		return fmt.Errorf("invalid collection: %w", err)
	}

	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal collection: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write collection file: %w", err)
	}

	return nil
}

func GetCollectionPath(fileStorage *storage.FileStorage, name string) string {
	filename := fmt.Sprintf("%s.json", name)
	return filepath.Join(fileStorage.CollectionsDir(), filename)
}

func validateCollection(collection *Collection) error {
	if collection.Name == "" {
		return fmt.Errorf("collection name is required")
	}

	if collection.Items == nil {
		collection.Items = make([]*types.RequestItem, 0)
	}

	for i, item := range collection.Items {
		if item.Name == "" {
			return fmt.Errorf("item %d: name is required", i)
		}
		if item.Request == nil {
			return fmt.Errorf("item %d: request is required", i)
		}
	}

	return nil
}
