package getman

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/KonnorFrik/getman/testutil/fixture"
	"github.com/KonnorFrik/getman/types"
)

func TestIntegrationPathPlaceholders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/42" {
			t.Errorf("expected path /users/42, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client, err := NewClient(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	client.SetGlobalVariable("baseUrl", server.URL)

	if err := client.SetPlaceholder("id", "42"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	execution, err := client.ExecuteRequest(&types.Request{
		Method: http.MethodGet,
		URL:    "{{baseUrl}}/users/:id",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if execution.Error != "" {
		t.Fatalf("unexpected execution error: %s", execution.Error)
	}

	if execution.Response == nil || execution.Response.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 response, got %+v", execution.Response)
	}
}

func TestIntegrationPathPlaceholderMissing(t *testing.T) {
	client, err := NewClient(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.ExecuteRequest(&types.Request{
		Method: http.MethodGet,
		URL:    "http://example.com/users/:id",
	})
	if !errors.Is(err, ErrPlaceholderNotFound) {
		t.Fatalf("expected ErrPlaceholderNotFound, got %v", err)
	}
}

func TestSaveCollectionRegistersPlaceholders(t *testing.T) {
	client, err := NewClient(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collection := fixture.CreateTestCollection("test", []*types.RequestItem{
		{
			Name: "Get user",
			Request: &types.Request{
				Method: http.MethodGet,
				URL:    "http://example.com/users/:id",
			},
		},
	})

	if err := client.SaveCollection(collection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	value, ok := client.GetPlaceholder("id")
	if !ok {
		t.Fatal("expected placeholder to be registered")
	}

	if value != "" {
		t.Errorf("expected empty placeholder value, got %q", value)
	}
}
