package vault

import (
	"fmt"
	"sort"
	"strings"
)

// Bookmark represents a saved Vault path with an optional label.
type Bookmark struct {
	Label string
	Path  string
	Mount string
}

// BookmarkResult holds the outcome of a bookmark operation.
type BookmarkResult struct {
	Bookmark Bookmark
	Added    bool
	Removed  bool
	Success  bool
	Err      error
}

// BookmarkStore holds a collection of bookmarks keyed by label.
type BookmarkStore struct {
	entries map[string]Bookmark
}

// NewBookmarkStore initializes an empty BookmarkStore.
func NewBookmarkStore() *BookmarkStore {
	return &BookmarkStore{entries: make(map[string]Bookmark)}
}

// Add inserts or replaces a bookmark by label.
func (s *BookmarkStore) Add(label, mount, path string) BookmarkResult {
	if label == "" || path == "" {
		return BookmarkResult{Success: false, Err: fmt.Errorf("label and path are required")}
	}
	b := Bookmark{Label: label, Path: path, Mount: mount}
	s.entries[label] = b
	return BookmarkResult{Bookmark: b, Added: true, Success: true}
}

// Remove deletes a bookmark by label.
func (s *BookmarkStore) Remove(label string) BookmarkResult {
	b, ok := s.entries[label]
	if !ok {
		return BookmarkResult{Success: false, Err: fmt.Errorf("bookmark %q not found", label)}
	}
	delete(s.entries, label)
	return BookmarkResult{Bookmark: b, Removed: true, Success: true}
}

// List returns all bookmarks sorted by label.
func (s *BookmarkStore) List() []Bookmark {
	out := make([]Bookmark, 0, len(s.entries))
	for _, b := range s.entries {
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Label) < strings.ToLower(out[j].Label)
	})
	return out
}

// Get retrieves a single bookmark by label.
func (s *BookmarkStore) Get(label string) (Bookmark, bool) {
	b, ok := s.entries[label]
	return b, ok
}
