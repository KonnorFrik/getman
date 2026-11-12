/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package environment

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/KonnorFrik/getman/errors"
)

var pathPlaceholderPattern = regexp.MustCompile(`(^|/):([A-Za-z_][A-Za-z0-9_]*)`)

// PlaceholderStore stores key-value pairs used to resolve path placeholders
// such as "/users/:userId".
type PlaceholderStore struct {
	mu     sync.RWMutex
	path   string
	values map[string]string
}

// NewPlaceholderStore creates a new placeholder store backed by path.
func NewPlaceholderStore(path string) *PlaceholderStore {
	return &PlaceholderStore{
		path:   path,
		values: make(map[string]string),
	}
}

// Load loads placeholders from disk. A missing file is treated as an empty store.
func (s *PlaceholderStore) Load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			s.mu.Lock()
			s.values = make(map[string]string)
			s.mu.Unlock()
			return nil
		}

		return fmt.Errorf("failed to read placeholders file: %w", err)
	}

	var values map[string]string
	if err := json.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("failed to parse placeholders file: %w", err)
	}

	if values == nil {
		values = make(map[string]string)
	}

	s.mu.Lock()
	s.values = values
	s.mu.Unlock()

	return nil
}

// Save writes placeholders to disk.
func (s *PlaceholderStore) Save() error {
	s.mu.RLock()
	data, err := json.MarshalIndent(s.values, "", "  ")
	s.mu.RUnlock()

	if err != nil {
		return fmt.Errorf("failed to marshal placeholders: %w", err)
	}

	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write placeholders file: %w", err)
	}

	return nil
}

// Set sets a placeholder value.
func (s *PlaceholderStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = value
}

// Ensure creates an empty placeholder if it does not exist.
func (s *PlaceholderStore) Ensure(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.values[key]; !ok {
		s.values[key] = ""
	}
}

// Get returns a placeholder value.
func (s *PlaceholderStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.values[key]
	return value, ok
}

// Delete removes a placeholder.
func (s *PlaceholderStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.values, key)
}

// Values returns a copy of all placeholders.
func (s *PlaceholderStore) Values() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]string, len(s.values))
	for key, value := range s.values {
		result[key] = value
	}

	return result
}

// Keys returns sorted placeholder names.
func (s *PlaceholderStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.values))
	for key := range s.values {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	return keys
}

// ExtractPathPlaceholders returns unique placeholder names found in a URL path.
func ExtractPathPlaceholders(rawURL string) []string {
	path := pathPart(rawURL)
	matches := pathPlaceholderPattern.FindAllStringSubmatch(path, -1)

	seen := make(map[string]struct{}, len(matches))
	result := make([]string, 0, len(matches))

	for _, match := range matches {
		name := match[2]
		if _, ok := seen[name]; ok {
			continue
		}

		seen[name] = struct{}{}
		result = append(result, name)
	}

	return result
}

// ResolvePath replaces path placeholders such as ":userId" with stored values.
func (s *PlaceholderStore) ResolvePath(rawURL string) (string, error) {
	path := pathPart(rawURL)
	suffix := rawURL[len(path):]

	missing := make(map[string]struct{})
	seenMissing := make([]string, 0)

	s.mu.RLock()
	resolved := pathPlaceholderPattern.ReplaceAllStringFunc(path, func(match string) string {
		submatch := pathPlaceholderPattern.FindStringSubmatch(match)
		if len(submatch) < 3 {
			return match
		}

		name := submatch[2]
		value, ok := s.values[name]

		if !ok || value == "" {
			if _, exists := missing[name]; !exists {
				missing[name] = struct{}{}
				seenMissing = append(seenMissing, name)
			}

			return match
		}

		return submatch[1] + value
	})
	s.mu.RUnlock()

	if len(seenMissing) > 0 {
		return "", fmt.Errorf("%w: %s", errors.ErrPlaceholderNotFound, strings.Join(seenMissing, ", "))
	}

	return resolved + suffix, nil
}

func pathPart(rawURL string) string {
	if idx := strings.IndexAny(rawURL, "?#"); idx >= 0 {
		return rawURL[:idx]
	}

	return rawURL
}
