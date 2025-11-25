package collections

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KonnorFrik/getman/storage"
	"github.com/KonnorFrik/getman/testutil/helper"
	"github.com/KonnorFrik/getman/types"
)

func TestUnitLoadCollectionFromFile_Valid(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	collectionJSON := `{
		"name": "Test Collection",
		"description": "Test description",
		"items": [
			{
				"name": "Test Request",
				"request": {
					"method": "GET",
					"url": "http://example.com"
				}
			}
		]
	}`

	filePath := filepath.Join(dir, "test.json")
	if err := os.WriteFile(filePath, []byte(collectionJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	collection, err := LoadCollectionFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if collection.Name != "Test Collection" {
		t.Errorf("expected name 'Test Collection', got %s", collection.Name)
	}

	if collection.Description != "Test description" {
		t.Errorf("expected description 'Test description', got %s", collection.Description)
	}

	if len(collection.Items) != 1 {
		t.Errorf("expected 1 item, got %d", len(collection.Items))
	}

	if collection.Items[0].Name != "Test Request" {
		t.Errorf("expected item name 'Test Request', got %s", collection.Items[0].Name)
	}
}

func TestUnitLoadCollectionFromFile_InvalidJSON(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	filePath := filepath.Join(dir, "test.json")
	if err := os.WriteFile(filePath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = LoadCollectionFromFile(filePath)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestUnitLoadCollectionFromFile_FileNotFound(t *testing.T) {
	_, err := LoadCollectionFromFile("/nonexistent/file.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestUnitLoadCollectionFromFile_InvalidCollection(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	collectionJSON := `{
		"name": "",
		"items": []
	}`

	filePath := filepath.Join(dir, "test.json")
	if err := os.WriteFile(filePath, []byte(collectionJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = LoadCollectionFromFile(filePath)
	if err == nil {
		t.Fatal("expected error for invalid collection")
	}
}

func TestUnitSaveCollectionToFile_Valid(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	collection := &Collection{
		Name:        "Test Collection",
		Description: "Test description",
		Items: []*types.RequestItem{
			{
				Name: "Test Request",
				Request: &types.Request{
					Method: "GET",
					URL:    "http://example.com",
				},
			},
		},
	}

	filePath := filepath.Join(dir, "test.json")
	if err := SaveCollectionToFile(collection, filePath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("expected file to be created")
	}

	loadedCollection, err := LoadCollectionFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error loading collection: %v", err)
	}

	if loadedCollection.Name != collection.Name {
		t.Errorf("expected name %s, got %s", collection.Name, loadedCollection.Name)
	}
}

func TestUnitSaveCollectionToFile_InvalidCollection(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	collection := &Collection{
		Name:  "",
		Items: []*types.RequestItem{},
	}

	filePath := filepath.Join(dir, "test.json")
	err = SaveCollectionToFile(collection, filePath)
	if err == nil {
		t.Fatal("expected error for invalid collection")
	}
}

func TestUnitValidateCollection_Valid(t *testing.T) {
	collection := &Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Test Request",
				Request: &types.Request{
					Method: "GET",
					URL:    "http://example.com",
				},
			},
		},
	}

	if err := validateCollection(collection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitValidateCollection_MissingName(t *testing.T) {
	collection := &Collection{
		Name:  "",
		Items: []*types.RequestItem{},
	}

	err := validateCollection(collection)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestUnitValidateCollection_MissingItemName(t *testing.T) {
	collection := &Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "",
				Request: &types.Request{
					Method: "GET",
					URL:    "http://example.com",
				},
			},
		},
	}

	err := validateCollection(collection)
	if err == nil {
		t.Fatal("expected error for missing item name")
	}
}

func TestUnitValidateCollection_MissingRequest(t *testing.T) {
	collection := &Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name:    "Test Request",
				Request: nil,
			},
		},
	}

	err := validateCollection(collection)
	if err == nil {
		t.Fatal("expected error for missing request")
	}
}

func TestUnitGetCollectionPath(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	fileStorage, err := storage.NewFileStorage(dir)
	if err != nil {
		t.Fatalf("failed to create file storage: %v", err)
	}

	path := GetCollectionPath(fileStorage, "test")
	expectedPath := filepath.Join(dir, "collections", "test.json")

	if path != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, path)
	}
}

func TestUnitValidateCollection_EmptyItems(t *testing.T) {
	collection := &Collection{
		Name:  "Test Collection",
		Items: nil,
	}

	if err := validateCollection(collection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if collection.Items == nil {
		t.Fatal("expected items to be initialized")
	}

	if len(collection.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(collection.Items))
	}
}

func TestUnitLoadCollectionFromFile_MultipleItems(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	collectionJSON := `{
		"name": "Test Collection",
		"items": [
			{
				"name": "Request 1",
				"request": {
					"method": "GET",
					"url": "http://example.com/1"
				}
			},
			{
				"name": "Request 2",
				"request": {
					"method": "POST",
					"url": "http://example.com/2"
				}
			}
		]
	}`

	filePath := filepath.Join(dir, "test.json")
	if err := os.WriteFile(filePath, []byte(collectionJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	collection, err := LoadCollectionFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(collection.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(collection.Items))
	}

	if collection.Items[0].Name != "Request 1" {
		t.Errorf("expected item 0 name 'Request 1', got %s", collection.Items[0].Name)
	}

	if collection.Items[1].Name != "Request 2" {
		t.Errorf("expected item 1 name 'Request 2', got %s", collection.Items[1].Name)
	}
}

func TestUnitSaveCollectionToFile_EmptyItems(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	collection := &Collection{
		Name:  "Test Collection",
		Items: []*types.RequestItem{},
	}

	filePath := filepath.Join(dir, "test.json")
	if err := SaveCollectionToFile(collection, filePath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loadedCollection, err := LoadCollectionFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error loading collection: %v", err)
	}

	if len(loadedCollection.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(loadedCollection.Items))
	}
}
