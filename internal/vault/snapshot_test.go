package vault

import (
	"errors"
	"testing"
	"time"
)

func TestSnapshotResult_Fields(t *testing.T) {
	now := time.Now()
	r := SnapshotResult{
		TakenAt: now,
		Success: true,
		Entries: []SnapshotEntry{
			{Path: "secret/foo", Data: map[string]interface{}{"key": "val"}},
		},
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if r.TakenAt != now {
		t.Error("expected TakenAt to match")
	}
	if len(r.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(r.Entries))
	}
	if r.Entries[0].Path != "secret/foo" {
		t.Errorf("unexpected path: %s", r.Entries[0].Path)
	}
}

func TestSnapshotResult_FailureState(t *testing.T) {
	r := SnapshotResult{
		Success: false,
		Err:     errors.New("walk failed: connection refused"),
	}
	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Err == nil {
		t.Error("expected Err to be set")
	}
}

func TestSnapshotEntry_EmptyData(t *testing.T) {
	e := SnapshotEntry{Path: "secret/empty", Data: map[string]interface{}{}}
	if len(e.Data) != 0 {
		t.Error("expected empty data map")
	}
}

func TestSnapshotEntry_NilData(t *testing.T) {
	e := SnapshotEntry{Path: "secret/nil"}
	if e.Data != nil {
		t.Error("expected nil data")
	}
}

func TestSnapshotResult_ZeroEntries(t *testing.T) {
	r := SnapshotResult{Success: true}
	if len(r.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(r.Entries))
	}
}
