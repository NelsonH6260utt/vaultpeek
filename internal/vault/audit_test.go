package vault

import (
	"testing"
	"time"
)

func baseAuditLog() *AuditLog {
	now := time.Now()
	return &AuditLog{
		Mount: "secret",
		Entries: []AuditEntry{
			{Path: "secret/foo", Operation: "read", Version: 1, Timestamp: now.Add(-2 * time.Hour)},
			{Path: "secret/bar", Operation: "write", Version: 2, Timestamp: now.Add(-1 * time.Hour)},
			{Path: "secret/foo", Operation: "write", Version: 3, Timestamp: now},
		},
	}
}

func TestAuditLog_Summary(t *testing.T) {
	log := baseAuditLog()
	s := log.Summary()
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
	if s != "mount=secret entries=3" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestAuditLog_FilterByPath(t *testing.T) {
	log := baseAuditLog()
	results := log.FilterByPath("secret/foo")
	if len(results) != 2 {
		t.Errorf("expected 2 entries for secret/foo, got %d", len(results))
	}
}

func TestAuditLog_FilterByPath_NoMatch(t *testing.T) {
	log := baseAuditLog()
	results := log.FilterByPath("secret/missing")
	if len(results) != 0 {
		t.Errorf("expected 0 entries, got %d", len(results))
	}
}

func TestAuditLog_FilterByOperation(t *testing.T) {
	log := baseAuditLog()
	results := log.FilterByOperation("write")
	if len(results) != 2 {
		t.Errorf("expected 2 write entries, got %d", len(results))
	}
}

func TestAuditLog_SortByTime(t *testing.T) {
	log := baseAuditLog()
	// Reverse order first
	log.Entries[0], log.Entries[2] = log.Entries[2], log.Entries[0]
	log.SortByTime()
	for i := 1; i < len(log.Entries); i++ {
		if log.Entries[i].Timestamp.Before(log.Entries[i-1].Timestamp) {
			t.Errorf("entries not sorted at index %d", i)
		}
	}
}

func TestAuditEntry_Fields(t *testing.T) {
	now := time.Now()
	e := AuditEntry{
		Path:      "secret/test",
		Operation: "delete",
		Version:   5,
		Timestamp: now,
		Renewable: false,
	}
	if e.Path != "secret/test" {
		t.Errorf("unexpected path: %s", e.Path)
	}
	if e.Version != 5 {
		t.Errorf("unexpected version: %d", e.Version)
	}
}
