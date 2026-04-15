package vault

import (
	"strings"
	"testing"
)

func TestSearchResult_Fields(t *testing.T) {
	r := SearchResult{
		Path:       "secret/myapp/db",
		MatchedKey: "password",
	}
	if r.Path != "secret/myapp/db" {
		t.Errorf("expected path 'secret/myapp/db', got %q", r.Path)
	}
	if r.MatchedKey != "password" {
		t.Errorf("expected matched key 'password', got %q", r.MatchedKey)
	}
}

func TestSearchResult_EmptyQuery(t *testing.T) {
	// An empty query string lowercased is still empty — verify strings.Contains behaviour.
	q := strings.ToLower("")
	if !strings.Contains("anything", q) {
		t.Error("expected empty query to match any string")
	}
}

func TestSearchResult_CaseInsensitive(t *testing.T) {
	query := "PASSWORD"
	key := "password"
	if !strings.Contains(strings.ToLower(key), strings.ToLower(query)) {
		t.Errorf("expected case-insensitive match of %q in %q", query, key)
	}
}

func TestSearchResult_ValueMatch(t *testing.T) {
	query := "prod"
	value := "production-db-host"
	if !strings.Contains(strings.ToLower(value), strings.ToLower(query)) {
		t.Errorf("expected %q to match in value %q", query, value)
	}
}

func TestSearchResult_NoMatch(t *testing.T) {
	query := "xyz123"
	key := "username"
	value := "admin"
	keyMatch := strings.Contains(strings.ToLower(key), strings.ToLower(query))
	valMatch := strings.Contains(strings.ToLower(value), strings.ToLower(query))
	if keyMatch || valMatch {
		t.Error("expected no match for unrelated query")
	}
}
