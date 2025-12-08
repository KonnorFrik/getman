/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package getman

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Storage  StorageConfig  `yaml:"storage"`
	Defaults DefaultsConfig `yaml:"defaults"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type StorageConfig struct {
	BasePath string `yaml:"base_path"`
}

type DefaultsConfig struct {
	Timeout TimeoutConfig `yaml:"timeout"`
	Cookies CookiesConfig `yaml:"cookies"`
}

type TimeoutConfig struct {
	Connect time.Duration `yaml:"connect"`
	Read    time.Duration `yaml:"read"`
}

type CookiesConfig struct {
	AutoManage bool `yaml:"auto_manage"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *Config, configPath string) error {
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
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
			Level:  "info",
			Format: "text",
		},
	}
}

func validateConfig(config *Config) error {
	if config.Storage.BasePath == "" {
		return fmt.Errorf("storage.base_path is required")
	}

	if config.Defaults.Timeout.Connect <= 0 {
		return fmt.Errorf("defaults.timeout.connect must be positive")
	}

	if config.Defaults.Timeout.Read <= 0 {
		return fmt.Errorf("defaults.timeout.read must be positive")
	}

	if config.Logging.Level == "" {
		return fmt.Errorf("logging.level is required")
	}

	if config.Logging.Format == "" {
		return fmt.Errorf("logging.format is required")
	}

	return nil
}
