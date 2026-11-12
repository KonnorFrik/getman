package environment

import (
	"errors"
	"path/filepath"
	"testing"

	getmanerrors "github.com/KonnorFrik/getman/errors"
)

func TestExtractPathPlaceholders(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want []string
	}{
		{
			name: "path placeholders",
			url:  "https://api.example.com/users/:id/posts/:postID",
			want: []string{"id", "postID"},
		},
		{
			name: "does not match scheme or port",
			url:  "http://localhost:8080/users/:id",
			want: []string{"id"},
		},
		{
			name: "ignores query and fragment",
			url:  "https://api.example.com/users/:id?next=/:other#frag",
			want: []string{"id"},
		},
		{
			name: "deduplicates",
			url:  "/users/:id/posts/:id",
			want: []string{"id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractPathPlaceholders(tt.url)
			if len(got) != len(tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("expected %v, got %v", tt.want, got)
				}
			}
		})
	}
}

func TestPlaceholderStore_ResolvePath(t *testing.T) {
	store := NewPlaceholderStore(filepath.Join(t.TempDir(), "placeholders.json"))
	store.Set("id", "42")

	got, err := store.ResolvePath("https://api.example.com/users/:id?verbose=true")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := "https://api.example.com/users/42?verbose=true"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestPlaceholderStore_ResolvePathMissing(t *testing.T) {
	store := NewPlaceholderStore(filepath.Join(t.TempDir(), "placeholders.json"))

	_, err := store.ResolvePath("https://api.example.com/users/:id")
	if !errors.Is(err, getmanerrors.ErrPlaceholderNotFound) {
		t.Fatalf("expected ErrPlaceholderNotFound, got %v", err)
	}
}

func TestPlaceholderStore_ResolvePathEmpty(t *testing.T) {
	store := NewPlaceholderStore(filepath.Join(t.TempDir(), "placeholders.json"))
	store.Ensure("id")

	_, err := store.ResolvePath("https://api.example.com/users/:id")
	if !errors.Is(err, getmanerrors.ErrPlaceholderNotFound) {
		t.Fatalf("expected ErrPlaceholderNotFound, got %v", err)
	}
}

func TestPlaceholderStore_SaveLoad(t *testing.T) {
	path := filepath.Join(t.TempDir(), "placeholders.json")
	store := NewPlaceholderStore(path)
	store.Set("id", "42")
	store.Set("slug", "hello")

	if err := store.Save(); err != nil {
		t.Fatalf("unexpected error saving: %v", err)
	}

	loaded := NewPlaceholderStore(path)
	if err := loaded.Load(); err != nil {
		t.Fatalf("unexpected error loading: %v", err)
	}

	if value, ok := loaded.Get("id"); !ok || value != "42" {
		t.Errorf("expected id=42, got %q (exists=%t)", value, ok)
	}

	if value, ok := loaded.Get("slug"); !ok || value != "hello" {
		t.Errorf("expected slug=hello, got %q (exists=%t)", value, ok)
	}
}
