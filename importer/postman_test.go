package importer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KonnorFrik/getman/testutil"
)

func TestUnitImportFromPostman_Valid(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := testutil.GetTestPostmanCollectionJSON()

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

func TestUnitImportFromPostman_InvalidJSON(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	filePath := filepath.Join(dir, "postman.json")
	if err := os.WriteFile(filePath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = ImportFromPostman(filePath)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestUnitImportFromPostman_FileNotFound(t *testing.T) {
	_, err := ImportFromPostman("/nonexistent/file.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestUnitImportFromPostman_WithHeaders(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Test Request",
				"request": {
					"method": "GET",
					"header": [
						{
							"key": "Accept",
							"value": "application/json"
						},
						{
							"key": "Authorization",
							"value": "Bearer token123"
						}
					],
					"url": {
						"raw": "http://example.com"
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

	req := collection.Items[0].Request
	if req.Headers["Accept"] != "application/json" {
		t.Errorf("expected Accept header 'application/json', got %s", req.Headers["Accept"])
	}

	if req.Headers["Authorization"] != "Bearer token123" {
		t.Errorf("expected Authorization header 'Bearer token123', got %s", req.Headers["Authorization"])
	}
}

func TestUnitImportFromPostman_WithBody_RAW(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Test Request",
				"request": {
					"method": "POST",
					"body": {
						"mode": "raw",
						"raw": "{\"key\": \"value\"}"
					},
					"url": {
						"raw": "http://example.com"
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

	req := collection.Items[0].Request
	if req.Body == nil {
		t.Fatal("expected body to be set")
	}

	if req.Body.Type != "raw" {
		t.Errorf("expected body type 'raw', got %s", req.Body.Type)
	}

	if string(req.Body.Content) != `{"key": "value"}` {
		t.Errorf("expected body content '{\"key\": \"value\"}', got %s", string(req.Body.Content))
	}
}

func TestUnitImportFromPostman_WithBody_FormData(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Test Request",
				"request": {
					"method": "POST",
					"body": {
						"mode": "formdata",
						"formdata": [
							{
								"key": "field1",
								"value": "value1"
							},
							{
								"key": "field2",
								"value": "value2"
							}
						]
					},
					"url": {
						"raw": "http://example.com"
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

	req := collection.Items[0].Request
	if req.Body == nil {
		t.Fatal("expected body to be set")
	}

	if req.Body.Type != "formdata" {
		t.Errorf("expected body type 'formdata', got %s", req.Body.Type)
	}
}

func TestUnitImportFromPostman_WithBody_URLEncoded(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Test Request",
				"request": {
					"method": "POST",
					"body": {
						"mode": "urlencoded",
						"urlencoded": [
							{
								"key": "field1",
								"value": "value1"
							}
						]
					},
					"url": {
						"raw": "http://example.com"
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

	req := collection.Items[0].Request
	if req.Body == nil {
		t.Fatal("expected body to be set")
	}

	if req.Body.Type != "urlencoded" {
		t.Errorf("expected body type 'urlencoded', got %s", req.Body.Type)
	}
}

func TestUnitImportFromPostman_WithBasicAuth(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Test Request",
				"request": {
					"method": "GET",
					"auth": {
						"type": "basic",
						"basic": [
							{
								"key": "username",
								"value": "testuser"
							},
							{
								"key": "password",
								"value": "testpass"
							}
						]
					},
					"url": {
						"raw": "http://example.com"
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

	req := collection.Items[0].Request
	if req.Auth == nil {
		t.Fatal("expected auth to be set")
	}

	if req.Auth.Type != "basic" {
		t.Errorf("expected auth type 'basic', got %s", req.Auth.Type)
	}

	if req.Auth.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", req.Auth.Username)
	}

	if req.Auth.Password != "testpass" {
		t.Errorf("expected password 'testpass', got %s", req.Auth.Password)
	}
}

func TestUnitImportFromPostman_WithBearerAuth(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Test Request",
				"request": {
					"method": "GET",
					"auth": {
						"type": "bearer",
						"bearer": [
							{
								"key": "token",
								"value": "testtoken123"
							}
						]
					},
					"url": {
						"raw": "http://example.com"
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

	req := collection.Items[0].Request
	if req.Auth == nil {
		t.Fatal("expected auth to be set")
	}

	if req.Auth.Type != "bearer" {
		t.Errorf("expected auth type 'bearer', got %s", req.Auth.Type)
	}

	if req.Auth.Token != "testtoken123" {
		t.Errorf("expected token 'testtoken123', got %s", req.Auth.Token)
	}
}

func TestUnitImportFromPostman_WithAPIKeyAuth(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Test Request",
				"request": {
					"method": "GET",
					"auth": {
						"type": "apikey",
						"apikey": [
							{
								"key": "key",
								"value": "X-API-Key"
							},
							{
								"key": "value",
								"value": "testapikey123"
							},
							{
								"key": "in",
								"value": "header"
							}
						]
					},
					"url": {
						"raw": "http://example.com"
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

	req := collection.Items[0].Request
	if req.Auth == nil {
		t.Fatal("expected auth to be set")
	}

	if req.Auth.Type != "apikey" {
		t.Errorf("expected auth type 'apikey', got %s", req.Auth.Type)
	}

	if req.Auth.APIKey != "testapikey123" {
		t.Errorf("expected API key 'testapikey123', got %s", req.Auth.APIKey)
	}

	if req.Auth.KeyName != "X-API-Key" {
		t.Errorf("expected key name 'X-API-Key', got %s", req.Auth.KeyName)
	}

	if req.Auth.Location != "header" {
		t.Errorf("expected location 'header', got %s", req.Auth.Location)
	}
}

func TestUnitImportFromPostman_NestedItems(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Folder",
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
							"method": "GET",
							"url": {
								"raw": "http://example.com/2"
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

func TestUnitImportFromPostman_APIKeyAuth_QueryLocation(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	postmanJSON := `{
		"info": {
			"name": "Test Collection",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Test Request",
				"request": {
					"method": "GET",
					"auth": {
						"type": "apikey",
						"apikey": [
							{
								"key": "key",
								"value": "api_key"
							},
							{
								"key": "value",
								"value": "testapikey123"
							},
							{
								"key": "in",
								"value": "query"
							}
						]
					},
					"url": {
						"raw": "http://example.com"
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

	req := collection.Items[0].Request
	if req.Auth.Location != "query" {
		t.Errorf("expected location 'query', got %s", req.Auth.Location)
	}
}

