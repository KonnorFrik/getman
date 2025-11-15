package environment

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"sync"
)

type Environment struct {
	mu sync.RWMutex
	Name      string            `json:"name"`
	Variables map[string]string `json:"variables"`
}

func NewEnvironment(name string) *Environment {
	return &Environment{
		Name: name,
	}
}

// NewEnvironmentFromFile - load env directly from file.
func NewEnvironmentFromFile(filepath string) (*Environment, error) {
	data, err := os.ReadFile(filepath)

	if err != nil {
		return nil, fmt.Errorf("failed to read environment file: %w", err)
	}

	var env Environment
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("failed to parse environment file: %w", err)
	}

	if err := validateEnvironment(&env); err != nil {
		return nil, fmt.Errorf("invalid environment: %w", err)
	}

	return &env, nil
}

// Set - save a pair key-value in environment 'e'.
// if key already exist - it will be overwriten.
func (e *Environment) Set(key, value string) { 
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Variables[key] = value
}

// Get - return a value binded with a 'key'.
// if key don't exist - reutrn zero-value and false.
func (e *Environment) Get(key string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	value, ok := e.Variables[key]
	return value, ok
}

// Clear - delete all environment 'e'.
func (e *Environment) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Variables = make(map[string]string)
}

func (e *Environment) Delete(key string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.Variables, key)
}

func (e *Environment) CopyMap() map[string]string {
	e.mu.Lock()
	defer e.mu.Unlock()
	result := make(map[string]string, len(e.Variables))
	maps.Copy(result, e.Variables)
	return result
}

// Load - overwrite 'e' with data from 'filepath'.
func (e *Environment) Load(filepath string) error {
	data, err := os.ReadFile(filepath)

	if err != nil {
		return fmt.Errorf("failed to read environment file: %w", err)
	}

	var env Environment
	if err := json.Unmarshal(data, &env); err != nil {
		return fmt.Errorf("failed to parse environment file: %w", err)
	}

	if err := validateEnvironment(&env); err != nil {
		return fmt.Errorf("invalid environment: %w", err)
	}

	e.Name = env.Name
	e.Variables = env.Variables
	return nil
}

// Save - save environment 'e' into file.
// if file exist - it will be overwriten.
func (e *Environment) Save(filepath string) error {
	if err := validateEnvironment(e); err != nil {
		return fmt.Errorf("invalid environment: %w", err)
	}

	data, err := json.MarshalIndent(e, "", "  ")

	if err != nil {
		return fmt.Errorf("failed to marshal environment: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write environment file: %w", err)
	}

	return nil
}

// func LoadEnvironmentFromFile(filePath string) (*types.Environment, error) {
// 	data, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read environment file: %w", err)
// 	}
//
// 	var env types.Environment
// 	if err := json.Unmarshal(data, &env); err != nil {
// 		return nil, fmt.Errorf("failed to parse environment file: %w", err)
// 	}
//
// 	if err := validateEnvironment(&env); err != nil {
// 		return nil, fmt.Errorf("invalid environment: %w", err)
// 	}
//
// 	return &env, nil
// }
//
// func SaveEnvironmentToFile(env *types.Environment, filePath string) error {
// 	if err := validateEnvironment(env); err != nil {
// 		return fmt.Errorf("invalid environment: %w", err)
// 	}
//
// 	data, err := json.MarshalIndent(env, "", "  ")
// 	if err != nil {
// 		return fmt.Errorf("failed to marshal environment: %w", err)
// 	}
//
// 	if err := os.WriteFile(filePath, data, 0644); err != nil {
// 		return fmt.Errorf("failed to write environment file: %w", err)
// 	}
//
// 	return nil
// }

func validateEnvironment(env *Environment) error {
	if env.Name == "" {
		return fmt.Errorf("environment name is required")
	}

	if env.Variables == nil {
		env.Variables = make(map[string]string)
	}

	return nil
}
