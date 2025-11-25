package core

import (
	"testing"

	"github.com/KonnorFrik/getman/environment"
)

func TestIntegrationResolve_InRequest(t *testing.T) {
	envG := environment.NewEnvironment("global")
	envL := environment.NewEnvironment("global")
	envG.Set("baseUrl", "http://example.com")
	envL.Set("token", "testtoken123")
	resolver, err := NewVariableResolver(envG, envL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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
	envG := environment.NewEnvironment("global")
	envL := environment.NewEnvironment("global")
	envG.Set("testVar", "globalValue")
	envL.Set("testVar", "envValue")
	resolver, err := NewVariableResolver(envG, envL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resolver.Resolve("{{testVar}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "envValue" {
		t.Errorf("expected 'envValue' (env priority), got %s", result)
	}
}

func TestIntegrationResolve_ComplexTemplate(t *testing.T) {
	envG := environment.NewEnvironment("global")
	envL := environment.NewEnvironment("global")
	envG.Set("protocol", "https")
	envG.Set("host", "api.example.com")
	envG.Set("path", "users")
	envL.Set("userId", "123")
	resolver, err := NewVariableResolver(envG, envL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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
	envG := environment.NewEnvironment("global")
	envL := environment.NewEnvironment("global")
	envG.Set("baseUrl", "http://example.com")
	envL.Set("apiPath", "/api/v1")
	envL.Set("resource", "users")
	resolver, err := NewVariableResolver(envG, envL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

