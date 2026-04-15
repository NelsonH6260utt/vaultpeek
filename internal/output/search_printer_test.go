package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/vaultpeek/internal/vault"
)

func TestPrintSearchResults_NoResults(t *testing.T) {
	var buf bytes.Buffer
	PrintSearchResults(&buf, nil, "myquery")
	out := buf.String()
	if !strings.Contains(out, "No results found") {
		t.Errorf("expected 'No results found', got: %s", out)
	}
	if !strings.Contains(out, "myquery") {
		t.Errorf("expected query 'myquery' in output, got: %s", out)
	}
}

func TestPrintSearchResults_ShowsPath(t *testing.T) {
	var buf bytes.Buffer
	results := []vault.SearchResult{
		{Path: "secret/app/config", MatchedKey: "db_password"},
	}
	PrintSearchResults(&buf, results, "password")
	out := buf.String()
	if !strings.Contains(out, "secret/app/config") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "db_password") {
		t.Errorf("expected matched key in output, got: %s", out)
	}
}

func TestPrintSearchResults_ShowsCount(t *testing.T) {
	var buf bytes.Buffer
	results := []vault.SearchResult{
		{Path: "secret/a", MatchedKey: "key1"},
		{Path: "secret/b", MatchedKey: "key2"},
	}
	PrintSearchResults(&buf, results, "key")
	out := buf.String()
	if !strings.Contains(out, "2 match") {
		t.Errorf("expected match count in output, got: %s", out)
	}
}

func TestPrintSearchResults_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	results := []vault.SearchResult{
		{Path: "secret/z", MatchedKey: "token"},
		{Path: "secret/a", MatchedKey: "token"},
	}
	PrintSearchResults(&buf, results, "token")
	out := buf.String()
	idxA := strings.Index(out, "secret/a")
	idxZ := strings.Index(out, "secret/z")
	if idxA > idxZ {
		t.Error("expected secret/a to appear before secret/z in sorted output")
	}
}
