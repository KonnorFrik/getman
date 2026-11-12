package importer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KonnorFrik/getman/collections"
	"github.com/KonnorFrik/getman/types"
)

func TestExportToPostman_RoundTrip(t *testing.T) {
	dir := t.TempDir()

	collection := &collections.Collection{
		Name:        "Export Test",
		Description: "round trip",
		Items: []*types.RequestItem{
			{
				Name: "Create User",
				Request: &types.Request{
					Method: "POST",
					URL:    "{{baseUrl}}/users",
					Headers: map[string]string{
						"Accept": "application/json",
					},
					Body: &types.RequestBody{
						Type:        "json",
						Content:     []byte(`{"name":"John"}`),
						ContentType: "application/json",
					},
					Auth: &types.Auth{Type: "bearer", Token: "{{token}}"},
				},
			},
		},
	}

	filePath := filepath.Join(dir, "export.json")
	if err := ExportToPostman(collection, filePath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected exported file: %v", err)
	}

	imported, err := ImportFromPostman(filePath)
	if err != nil {
		t.Fatalf("unexpected import error: %v", err)
	}

	if imported.Name != collection.Name {
		t.Errorf("expected name %q, got %q", collection.Name, imported.Name)
	}

	if len(imported.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(imported.Items))
	}

	req := imported.Items[0].Request
	if req.URL != "{{baseUrl}}/users" {
		t.Errorf("expected URL to round trip, got %q", req.URL)
	}

	if req.Auth == nil || req.Auth.Token != "{{token}}" {
		t.Errorf("expected bearer auth to round trip, got %+v", req.Auth)
	}

	if req.Body == nil || string(req.Body.Content) != `{"name":"John"}` {
		t.Errorf("expected body to round trip, got %+v", req.Body)
	}
}

func TestImportFromPostman_URLFromComponents(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "postman.json")

	data := `{
		"info": {"name": "URL Test"},
		"item": [{
			"name": "Get User",
			"request": {
				"method": "GET",
				"url": {
					"protocol": "https",
					"host": ["api", "example", "com"],
					"path": ["users", "42"],
					"query": [{"key": "verbose", "value": "true"}]
				}
			}
		}]
	}`

	if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	collection, err := ImportFromPostman(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "https://api.example.com/users/42?verbose=true"
	if got := collection.Items[0].Request.URL; got != expected {
		t.Errorf("expected URL %q, got %q", expected, got)
	}
}

func TestExportToPostman_InvalidCollection(t *testing.T) {
	if err := ExportToPostman(nil, filepath.Join(t.TempDir(), "out.json")); err == nil {
		t.Fatal("expected error for nil collection")
	}
}

func TestImportFromPostman_PathPlaceholderFromPathBlock(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "postman.json")

	data := `{
		"info": {"name": "Path Placeholder Test"},
		"item": [{
			"name": "Get Order",
			"request": {
				"method": "GET",
				"url": {
					"raw": "http://example.com/orders",
					"protocol": "http",
					"host": ["example", "com"],
					"path": ["orders", ":orderId"]
				}
			}
		}]
	}`

	if err := os.WriteFile(filePath, []byte(data), 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	collection, err := ImportFromPostman(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "http://example.com/orders/:orderId"
	if got := collection.Items[0].Request.URL; got != expected {
		t.Errorf("expected URL %q, got %q", expected, got)
	}
}
