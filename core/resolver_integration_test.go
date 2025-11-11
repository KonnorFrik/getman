package core

import (
	"testing"

	"github.com/KonnorFrik/getman/variables"
)

func TestIntegrationResolve_InRequest(t *testing.T) {
	store := variables.NewVariableStore()
	store.SetGlobal("baseUrl", "http://example.com")
	store.SetEnv("token", "testtoken123")
	resolver := NewVariableResolver(store)

	url := "{{baseUrl}}/api/users"
	resolvedURL, err := resolver.Resolve(url)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedURL != "http://example.com/api/users" {
		t.Errorf("expected 'http://example.com/api/users', got %s", resolvedURL)
	}

	headers := map[string]string{
		"Authorization": "Bearer {{token}}",
	}

	resolvedHeaders, err := resolver.ResolveMap(headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedHeaders["Authorization"] != "Bearer testtoken123" {
		t.Errorf("expected 'Bearer testtoken123', got %s", resolvedHeaders["Authorization"])
	}
}

func TestIntegrationResolve_EnvPriority(t *testing.T) {
	store := variables.NewVariableStore()
	store.SetGlobal("testVar", "globalValue")
	store.SetEnv("testVar", "envValue")
	resolver := NewVariableResolver(store)

	result, err := resolver.Resolve("{{testVar}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "envValue" {
		t.Errorf("expected 'envValue' (env priority), got %s", result)
	}
}

func TestIntegrationResolve_ComplexTemplate(t *testing.T) {
	store := variables.NewVariableStore()
	store.SetGlobal("protocol", "https")
	store.SetGlobal("host", "api.example.com")
	store.SetGlobal("path", "users")
	store.SetEnv("userId", "123")
	resolver := NewVariableResolver(store)

	template := "{{protocol}}://{{host}}/{{path}}/{{userId}}"
	result, err := resolver.Resolve(template)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "https://api.example.com/users/123"
	if result != expected {
		t.Errorf("expected '%s', got %s", expected, result)
	}
}

func TestIntegrationResolve_NestedVariables(t *testing.T) {
	store := variables.NewVariableStore()
	store.SetGlobal("baseUrl", "http://example.com")
	store.SetEnv("apiPath", "/api/v1")
	store.SetEnv("resource", "users")
	resolver := NewVariableResolver(store)

	template := "{{baseUrl}}{{apiPath}}/{{resource}}"
	result, err := resolver.Resolve(template)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "http://example.com/api/v1/users"
	if result != expected {
		t.Errorf("expected '%s', got %s", expected, result)
	}
}

