package importer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KonnorFrik/getman/testutil/helper"
	"github.com/KonnorFrik/getman/testutil/fixture"
	"github.com/KonnorFrik/getman/testutil/http_server"
)

func TestIntegrationImportFromPostman_RealFile(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	postmanJSON := fixture.GetTestPostmanCollectionJSON()

	filePath := filepath.Join(dir, "postman.json")
	if err := os.WriteFile(filePath, []byte(postmanJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	collection, err := ImportFromPostman(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if collection.Name != "Test Collection" {
		t.Errorf("expected name 'Test Collection', got %s", collection.Name)
	}

	if len(collection.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(collection.Items))
	}
}

func TestIntegrationImportFromPostman_ExecuteImported(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Get Health",
				"request": {
					"method": "GET",
					"url": {
						"raw": "` + http_server.GetServerURL() + `/health"
					}
				}
			}
		]
	}`

	filePath := filepath.Join(dir, "postman.json")
	if err := os.WriteFile(filePath, []byte(postmanJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	collection, err := ImportFromPostman(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(collection.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(collection.Items))
	}

	if collection.Items[0].Request.URL != http_server.GetServerURL()+"/health" {
		t.Errorf("expected URL %s/health, got %s", http_server.GetServerURL(), collection.Items[0].Request.URL)
	}
}

func TestIntegrationImportFromPostman_ComplexCollection(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Complex Collection",
			"description": "Complex collection with nested items",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Folder 1",
				"item": [
					{
						"name": "Request 1",
						"request": {
							"method": "GET",
							"url": {
								"raw": "http://example.com/1"
							}
						}
					},
					{
						"name": "Request 2",
						"request": {
							"method": "POST",
							"url": {
								"raw": "http://example.com/2"
							},
							"body": {
								"mode": "raw",
								"raw": "{\"key\": \"value\"}"
							}
						}
					}
				]
			},
			{
				"name": "Folder 2",
				"item": [
					{
						"name": "Request 3",
						"request": {
							"method": "PUT",
							"url": {
								"raw": "http://example.com/3"
							},
							"auth": {
								"type": "bearer",
								"bearer": [
									{
										"key": "token",
										"value": "testtoken"
									}
								]
							}
						}
					}
				]
			}
		]
	}`

	filePath := filepath.Join(dir, "postman.json")
	if err := os.WriteFile(filePath, []byte(postmanJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	collection, err := ImportFromPostman(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if collection.Name != "Complex Collection" {
		t.Errorf("expected name 'Complex Collection', got %s", collection.Name)
	}

	if len(collection.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(collection.Items))
	}

	if collection.Items[0].Name != "Request 1" {
		t.Errorf("expected item 0 name 'Request 1', got %s", collection.Items[0].Name)
	}

	if collection.Items[1].Name != "Request 2" {
		t.Errorf("expected item 1 name 'Request 2', got %s", collection.Items[1].Name)
	}

	if collection.Items[2].Name != "Request 3" {
		t.Errorf("expected item 2 name 'Request 3', got %s", collection.Items[2].Name)
	}
}

