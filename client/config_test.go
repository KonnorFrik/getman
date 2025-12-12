package getman

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/testutil/helper"
	"github.com/KonnorFrik/getman/testutil/fixture"
)

func TestUnitLoadConfig_Valid(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	configYAML := fixture.GetTestConfigYAML()

	filePath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(filePath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	config, err := LoadConfig(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.Storage.BasePath == "" {
		t.Error("expected base path to be set")
	}
}

func TestUnitLoadConfig_InvalidYAML(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	filePath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(filePath, []byte("invalid: yaml: ["), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = LoadConfig(filePath)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestUnitLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestUnitLoadConfig_InvalidConfig(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	configYAML := `storage:
  base_path: ""

defaults:
  timeout:
    connect: 0s
    read: 0s
  cookies:
    auto_manage: true

logging:
  level: ""
  format: ""
`

	filePath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(filePath, []byte(configYAML), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err = LoadConfig(filePath)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestUnitSaveConfig_Valid(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	config := DefaultConfig()

	filePath := filepath.Join(dir, "config.yaml")
	if err := SaveConfig(config, filePath); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("expected file to be created")
	}

	loadedConfig, err := LoadConfig(filePath)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	if loadedConfig.Storage.BasePath != config.Storage.BasePath {
		t.Errorf("expected base path %s, got %s", config.Storage.BasePath, loadedConfig.Storage.BasePath)
	}
}

func TestUnitSaveConfig_InvalidConfig(t *testing.T) {
	dir, err := helper.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer helper.CleanupTempDir(dir)

	config := &Config{
		Storage: StorageConfig{
			BasePath: "",
		},
		Defaults: DefaultsConfig{
			Timeout: TimeoutConfig{
				Connect: 0,
				Read:    0,
			},
		},
		Logging: LoggingConfig{
			Level:  "",
			Format: "",
		},
	}

	filePath := filepath.Join(dir, "config.yaml")
	err = SaveConfig(config, filePath)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestUnitDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Storage.BasePath == "" {
		t.Error("expected base path to be set")
	}

	if config.Defaults.Timeout.Connect <= 0 {
		t.Error("expected connect timeout to be positive")
	}

	if config.Defaults.Timeout.Read <= 0 {
		t.Error("expected read timeout to be positive")
	}

	if config.Defaults.Cookies.AutoManage != true {
		t.Error("expected auto manage cookies to be true")
	}

	if config.Logging.Level == "" {
		t.Error("expected logging level to be set")
	}

	if config.Logging.Format == "" {
		t.Error("expected logging format to be set")
	}
}

func TestUnitValidateConfig_Valid(t *testing.T) {
	config := DefaultConfig()

	if err := validateConfig(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitValidateConfig_MissingBasePath(t *testing.T) {
	config := &Config{
		Storage: StorageConfig{
			BasePath: "",
		},
		Defaults: DefaultsConfig{
			Timeout: TimeoutConfig{
				Connect: 30 * time.Second,
				Read:    30 * time.Second,
			},
			Cookies: CookiesConfig{
				AutoManage: true,
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Fatal("expected error for missing base path")
	}
}

func TestUnitValidateConfig_InvalidTimeout(t *testing.T) {
	config := &Config{
		Storage: StorageConfig{
			BasePath: "~/.getman",
		},
		Defaults: DefaultsConfig{
			Timeout: TimeoutConfig{
				Connect: 0,
				Read:    30 * time.Second,
			},
			Cookies: CookiesConfig{
				AutoManage: true,
			},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Fatal("expected error for invalid timeout")
	}
}

func TestUnitValidateConfig_MissingLoggingLevel(t *testing.T) {
	config := &Config{
		Storage: StorageConfig{
			BasePath: "~/.getman",
		},
		Defaults: DefaultsConfig{
			Timeout: TimeoutConfig{
				Connect: 30 * time.Second,
				Read:    30 * time.Second,
			},
			Cookies: CookiesConfig{
				AutoManage: true,
			},
		},
		Logging: LoggingConfig{
			Level:  "",
			Format: "text",
		},
	}

	err := validateConfig(config)
	if err == nil {
		t.Fatal("expected error for missing logging level")
	}
}

