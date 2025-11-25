package core

import (
	"testing"

	stderrors "errors"

	"github.com/KonnorFrik/getman/environment"
	"github.com/KonnorFrik/getman/errors"
)

func TestUnitResolve_SimpleVariable(t *testing.T) {
	envG := environment.NewEnvironment("global")
	envG.Set("testVar", "testValue")
	resolver, err := NewVariableResolver(envG, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resolver.Resolve("{{testVar}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "testValue" {
		t.Errorf("expected 'testValue', got %s", result)
	}
}

func TestUnitResolve_MultipleVariables(t *testing.T) {
	envG := environment.NewEnvironment("global")
	envG.Set("var1", "value1")
	envG.Set("var2", "value2")
	resolver, err := NewVariableResolver(envG, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resolver.Resolve("{{var1}} and {{var2}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "value1 and value2" {
		t.Errorf("expected 'value1 and value2', got %s", result)
	}
}

func TestUnitResolve_NoVariables(t *testing.T) {
	env := environment.NewEnvironment("global")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resolver.Resolve("no variables here")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "no variables here" {
		t.Errorf("expected 'no variables here', got %s", result)
	}
}

func TestUnitResolve_VariableNotFound(t *testing.T) {
	env := environment.NewEnvironment("global")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = resolver.Resolve("{{nonexistent}}")
	if err == nil {
		t.Fatal("expected error for nonexistent variable")
	}

	if !stderrors.Is(err, errors.ErrVariableNotFound) {
		t.Errorf("expected ErrVariableNotFound, got %v", err)
	}
}

func TestUnitResolve_EnvVariablePriority(t *testing.T) {
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

func TestUnitResolve_GlobalVariable(t *testing.T) {
	envG := environment.NewEnvironment("global")
	envG.Set("testVar", "globalValue")
	resolver, err := NewVariableResolver(envG, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resolver.Resolve("{{testVar}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "globalValue" {
		t.Errorf("expected 'globalValue', got %s", result)
	}
}

func TestUnitResolve_ComplexTemplate(t *testing.T) {
	envG := environment.NewEnvironment("global")
	envG.Set("baseUrl", "http://example.com")
	envG.Set("path", "api/users")
	envG.Set("id", "123")
	resolver, err := NewVariableResolver(envG, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resolver.Resolve("{{baseUrl}}/{{path}}/{{id}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "http://example.com/api/users/123"
	if result != expected {
		t.Errorf("expected '%s', got %s", expected, result)
	}
}

func TestUnitResolveMap(t *testing.T) {
	envG := environment.NewEnvironment("global")
	envG.Set("key", "testKey")
	envG.Set("value", "testValue")
	resolver, err := NewVariableResolver(envG, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]string{
		"{{key}}": "{{value}}",
	}

	result, err := resolver.ResolveMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["testKey"] != "testValue" {
		t.Errorf("expected map['testKey']='testValue', got %v", result)
	}
}

func TestUnitResolveMap_WithVariables(t *testing.T) {
	envG := environment.NewEnvironment("global")
	envG.Set("header", "Authorization")
	envG.Set("token", "Bearer abc123")
	resolver, err := NewVariableResolver(envG, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]string{
		"{{header}}": "{{token}}",
	}

	result, err := resolver.ResolveMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["Authorization"] != "Bearer abc123" {
		t.Errorf("expected map['Authorization']='Bearer abc123', got %v", result)
	}
}

func TestUnitResolveMap_VariableNotFound(t *testing.T) {
	env := environment.NewEnvironment("global")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]string{
		"key": "{{nonexistent}}",
	}

	_, err = resolver.ResolveMap(input)
	if err == nil {
		t.Fatal("expected error for nonexistent variable")
	}
}

func TestUnitValidateVariables(t *testing.T) {
	env := environment.NewEnvironment("global")
	env.Set("testVar", "testValue")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = resolver.ValidateVariables("{{testVar}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitValidateVariables_NotFound(t *testing.T) {
	env := environment.NewEnvironment("global")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = resolver.ValidateVariables("{{nonexistent}}")
	if err == nil {
		t.Fatal("expected error for nonexistent variable")
	}

	if !stderrors.Is(err, errors.ErrVariableNotFound) {
		t.Errorf("expected ErrVariableNotFound, got %v", err)
	}
}

func TestUnitValidateVariables_NoVariables(t *testing.T) {
	env := environment.NewEnvironment("global")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = resolver.ValidateVariables("no variables")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitValidateVariablesInMap(t *testing.T) {
	env := environment.NewEnvironment("global")
	env.Set("key", "testKey")
	env.Set("value", "testValue")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]string{
		"{{key}}": "{{value}}",
	}

	err = resolver.ValidateVariablesInMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitValidateVariablesInMap_NotFound(t *testing.T) {
	env := environment.NewEnvironment("global")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := map[string]string{
		"key": "{{nonexistent}}",
	}

	err = resolver.ValidateVariablesInMap(input)
	if err == nil {
		t.Fatal("expected error for nonexistent variable")
	}
}

func TestUnitResolve_EmptyVariable(t *testing.T) {
	env := environment.NewEnvironment("global")
	env.Set("emptyVar", "")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resolver.Resolve("{{emptyVar}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestUnitResolve_SameVariableMultipleTimes(t *testing.T) {
	env := environment.NewEnvironment("global")
	env.Set("var", "value")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resolver.Resolve("{{var}}-{{var}}-{{var}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "value-value-value" {
		t.Errorf("expected 'value-value-value', got %s", result)
	}
}

func TestUnitResolve_VariableWithSpaces(t *testing.T) {
	env := environment.NewEnvironment("global")
	env.Set("testVar", "testValue")
	resolver, err := NewVariableResolver(env, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := resolver.Resolve("{{ testVar }}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "testValue" {
		t.Errorf("expected 'testValue', got %s", result)
	}
}
