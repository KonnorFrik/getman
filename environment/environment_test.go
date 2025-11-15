package environment

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KonnorFrik/getman/testutil"
)

func TestUnitLoadEnvironmentFromFile_Valid(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	envJSON := `{
		"name": "test",
		"variables": {
			"key1": "value1",
			"key2": "value2"
		}
	}`

	filePath := filepath.Join(dir, "test.json")
	if err := os.WriteFile(filePath, []byte(envJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	env, err := NewEnvironmentFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env.Name != "test" {
		t.Errorf("expected name 'test', got %s", env.Name)
	}

	if env.Variables["key1"] != "value1" {
		t.Errorf("expected key1 'value1', got %s", env.Variables["key1"])
	}

	if env.Variables["key2"] != "value2" {
		t.Errorf("expected key2 'value2', got %s", env.Variables["key2"])
	}
}

func TestUnitLoadEnvironmentFromFile_InvalidJSON(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	filePath := filepath.Join(dir, "test.json")
	if err := os.WriteFile(filePath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = NewEnvironmentFromFile(filePath)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestUnitLoadEnvironmentFromFile_FileNotFound(t *testing.T) {
	_, err := NewEnvironmentFromFile("/nonexistent/file.json")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestUnitLoadEnvironmentFromFile_InvalidEnvironment(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	envJSON := `{
		"name": "",
		"variables": {}
	}`

	filePath := filepath.Join(dir, "test.json")
	if err := os.WriteFile(filePath, []byte(envJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = NewEnvironmentFromFile(filePath)
	if err == nil {
		t.Fatal("expected error for invalid environment")
	}
}

func TestUnitSaveEnvironmentToFile_Valid(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	env := &Environment{
		Name: "test",
		Variables: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	filePath := filepath.Join(dir, "test.json")
	if err := env.Save(filePath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("expected file to be created")
	}

	loadedEnv, err := NewEnvironmentFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error loading environment: %v", err)
	}

	if loadedEnv.Name != env.Name {
		t.Errorf("expected name %s, got %s", env.Name, loadedEnv.Name)
	}

	if loadedEnv.Variables["key1"] != env.Variables["key1"] {
		t.Errorf("expected key1 %s, got %s", env.Variables["key1"], loadedEnv.Variables["key1"])
	}
}

func TestUnitSaveEnvironmentToFile_InvalidEnvironment(t *testing.T) {
	dir, err := testutil.CreateTempDir()

	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	defer testutil.CleanupTempDir(dir)

	env := &Environment{
		Name:      "",
		Variables: map[string]string{},
	}

	filePath := filepath.Join(dir, "test.json")
	err = env.Save(filePath)

	if err == nil {
		t.Fatal("expected error for invalid environment")
	}
}

func TestUnitValidateEnvironment_Valid(t *testing.T) {
	env := &Environment{
		Name: "test",
		Variables: map[string]string{
			"key1": "value1",
		},
	}

	if err := validateEnvironment(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitValidateEnvironment_MissingName(t *testing.T) {
	env := &Environment{
		Name:      "",
		Variables: map[string]string{},
	}

	err := validateEnvironment(env)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestUnitValidateEnvironment_NilVariables(t *testing.T) {
	env := &Environment{
		Name:      "test",
		Variables: nil,
	}

	if err := validateEnvironment(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env.Variables == nil {
		t.Fatal("expected variables to be initialized")
	}
}

func TestUnitLoadEnvironmentFromFile_EmptyVariables(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	envJSON := `{
		"name": "test",
		"variables": {}
	}`

	filePath := filepath.Join(dir, "test.json")
	if err := os.WriteFile(filePath, []byte(envJSON), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	env, err := NewEnvironmentFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env.Name != "test" {
		t.Errorf("expected name 'test', got %s", env.Name)
	}

	if len(env.Variables) != 0 {
		t.Errorf("expected 0 variables, got %d", len(env.Variables))
	}
}

func TestUnitSaveEnvironmentToFile_EmptyVariables(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	env := &Environment{
		Name:      "test",
		Variables: map[string]string{},
	}

	filePath := filepath.Join(dir, "test.json")
	if err := env.Save(filePath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loadedEnv, err := NewEnvironmentFromFile(filePath)
	if err != nil {
		t.Fatalf("unexpected error loading environment: %v", err)
	}

	if len(loadedEnv.Variables) != 0 {
		t.Errorf("expected 0 variables, got %d", len(loadedEnv.Variables))
	}
}

