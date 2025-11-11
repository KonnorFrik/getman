package variables

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/KonnorFrik/getman/types"
)

func LoadEnvironmentFromFile(filePath string) (*types.Environment, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment file: %w", err)
	}

	var env types.Environment
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("failed to parse environment file: %w", err)
	}

	if err := validateEnvironment(&env); err != nil {
		return nil, fmt.Errorf("invalid environment: %w", err)
	}

	return &env, nil
}

func SaveEnvironmentToFile(env *types.Environment, filePath string) error {
	if err := validateEnvironment(env); err != nil {
		return fmt.Errorf("invalid environment: %w", err)
	}

	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal environment: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write environment file: %w", err)
	}

	return nil
}

func validateEnvironment(env *types.Environment) error {
	if env.Name == "" {
		return fmt.Errorf("environment name is required")
	}

	if env.Variables == nil {
		env.Variables = make(map[string]string)
	}

	return nil
}
