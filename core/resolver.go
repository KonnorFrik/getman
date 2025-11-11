package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/KonnorFrik/getman/errors"
	"github.com/KonnorFrik/getman/variables"
)

type VariableResolver struct {
	store *variables.VariableStore
}

func NewVariableResolver(store *variables.VariableStore) *VariableResolver {
	return &VariableResolver{
		store: store,
	}
}

var variablePattern = regexp.MustCompile(`\{\{([^}]+)\}\}`)

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
		value, ok := vr.store.Get(varName)
		if !ok {
			return "", fmt.Errorf("%w: %s", errors.ErrVariableNotFound, varName)
		}

		result = strings.ReplaceAll(result, match[0], value)
	}

	return result, nil
}

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
		if _, ok := vr.store.Get(varName); !ok {
			return fmt.Errorf("%w: %s", errors.ErrVariableNotFound, varName)
		}
	}

	return nil
}

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
