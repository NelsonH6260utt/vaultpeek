package vault

import (
	"fmt"
	"strings"
)

// AliasEntry represents a named shortcut to a Vault secret path.
type AliasEntry struct {
	Name string
	Path string
}

// AliasStore holds a collection of named path aliases.
type AliasStore struct {
	entries map[string]string
}

// AliasResult is returned by alias operations.
type AliasResult struct {
	Name    string
	Path    string
	Success bool
	Err     error
}

// NewAliasStore creates an empty AliasStore.
func NewAliasStore() *AliasStore {
	return &AliasStore{entries: make(map[string]string)}
}

// Set adds or updates an alias.
func (s *AliasStore) Set(name, path string) AliasResult {
	name = strings.TrimSpace(name)
	path = strings.TrimSpace(path)
	if name == "" {
		return AliasResult{Err: fmt.Errorf("alias name must not be empty")}
	}
	if path == "" {
		return AliasResult{Err: fmt.Errorf("alias path must not be empty")}
	}
	s.entries[name] = path
	return AliasResult{Name: name, Path: path, Success: true}
}

// Remove deletes an alias by name.
func (s *AliasStore) Remove(name string) AliasResult {
	if _, ok := s.entries[name]; !ok {
		return AliasResult{Name: name, Err: fmt.Errorf("alias %q not found", name)}
	}
	path := s.entries[name]
	delete(s.entries, name)
	return AliasResult{Name: name, Path: path, Success: true}
}

// Resolve returns the path for a given alias name.
func (s *AliasStore) Resolve(name string) (string, bool) {
	path, ok := s.entries[name]
	return path, ok
}

// List returns all aliases sorted by name.
func (s *AliasStore) List() []AliasEntry {
	keys := make([]string, 0, len(s.entries))
	for k := range s.entries {
		keys = append(keys, k)
	}
	sortStrings(keys)
	out := make([]AliasEntry, 0, len(keys))
	for _, k := range keys {
		out = append(out, AliasEntry{Name: k, Path: s.entries[k]})
	}
	return out
}

func sortStrings(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
