/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/KonnorFrik/getman/environment"
	"github.com/KonnorFrik/getman/errors"
)

// VariableResolver resolves variables from global and local environments.
type VariableResolver struct {
	global *environment.Environment
	local  *environment.Environment
}

// NewVariableResolver creates a new VariableResolver with global and optional local environment.
func NewVariableResolver(global, local *environment.Environment) (*VariableResolver, error) {
	if global == nil {
		return nil, fmt.Errorf("%w: global env can't be nil", errors.ErrInvalidArgument)
	}

	return &VariableResolver{
		local:  local,
		global: global,
	}, nil
}

var variablePattern = regexp.MustCompile(`\{\{([^}]+)\}\}`)

// Resolve resolves variables in a template string using the current environment.
func (vr *VariableResolver) Resolve(template string) (string, error) {
	matches := variablePattern.FindAllStringSubmatch(template, -1)

	if matches == nil {
		return template, nil
	}

	result := template

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		varName := strings.TrimSpace(match[1])
		var (
			value string
			ok    bool
		)

		if vr.local != nil {
			value, ok = vr.local.Get(varName)

		} else {
			value, ok = vr.global.Get(varName)
		}

		if !ok {
			return "", fmt.Errorf("%w: %s", errors.ErrVariableNotFound, varName)
		}

		result = strings.ReplaceAll(result, match[0], value)
	}

	return result, nil
}

// ResolveMap resolves variables in both keys and values of a map.
func (vr *VariableResolver) ResolveMap(m map[string]string) (map[string]string, error) {
	result := make(map[string]string)

	for k, v := range m {
		resolvedKey, err := vr.Resolve(k)

		if err != nil {
			return nil, err
		}

		resolvedValue, err := vr.Resolve(v)

		if err != nil {
			return nil, err
		}

		result[resolvedKey] = resolvedValue
	}

	return result, nil
}

// SetLocal sets the local environment for variable resolution.
func (vr *VariableResolver) SetLocal(local *environment.Environment) {
	vr.local = local
}

// SetGlobal sets the global environment for variable resolution.
func (vr *VariableResolver) SetGlobal(global *environment.Environment) {
	vr.global = global
}

// GetLocal returns the current local environment.
func (vr *VariableResolver) GetLocal() *environment.Environment {
	return vr.local
}

// GetGlobal returns the current global environment.
func (vr *VariableResolver) GetGlobal() *environment.Environment {
	return vr.global
}

// ValidateVariables validates that all variables in a template exist in the current environment.
func (vr *VariableResolver) ValidateVariables(template string) error {
	matches := variablePattern.FindAllStringSubmatch(template, -1)

	if matches == nil {
		return nil
	}

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		varName := strings.TrimSpace(match[1])

		var ok bool

		if vr.local != nil {
			_, ok = vr.local.Get(varName)

		} else {
			_, ok = vr.global.Get(varName)
		}

		if !ok {
			return fmt.Errorf("%w: %s", errors.ErrVariableNotFound, varName)
		}
	}

	return nil
}

// ValidateVariablesInMap validates variables in both keys and values of a map.
func (vr *VariableResolver) ValidateVariablesInMap(m map[string]string) error {
	for k, v := range m {
		if err := vr.ValidateVariables(k); err != nil {
			return err
		}

		if err := vr.ValidateVariables(v); err != nil {
			return err
		}
	}

	return nil
}
