package vault

import (
	"testing"
)

func TestPinResult_Fields(t *testing.T) {
	r := PinSecret("secret/data/app/db", 3)
	if !r.Success {
		t.Fatalf("expected success, got error: %v", r.Error)
	}
	if r.Path != "secret/data/app/db" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if r.Version != 3 {
		t.Errorf("expected version 3, got %d", r.Version)
	}
	if !r.Pinned {
		t.Error("expected Pinned to be true")
	}
	if r.PinnedAt.IsZero() {
		t.Error("expected PinnedAt to be set")
	}
}

func TestPinResult_FailureState(t *testing.T) {
	r := PinSecret("", 1)
	if r.Success {
		t.Error("expected failure for empty path")
	}
	if r.Error == nil {
		t.Error("expected non-nil error")
	}
}

func TestPinResult_InvalidVersion(t *testing.T) {
	r := PinSecret("secret/data/app", 0)
	if r.Success {
		t.Error("expected failure for version 0")
	}
	if r.Error == nil {
		t.Error("expected non-nil error for version 0")
	}

	r2 := PinSecret("secret/data/app", -5)
	if r2.Success {
		t.Error("expected failure for negative version")
	}
}

func TestIsPinned_True(t *testing.T) {
	meta := map[string]interface{}{"pinned_version": "2"}
	if !IsPinned(meta) {
		t.Error("expected IsPinned to return true")
	}
}

func TestIsPinned_False(t *testing.T) {
	meta := map[string]interface{}{"owner": "team-a"}
	if IsPinned(meta) {
		t.Error("expected IsPinned to return false")
	}
}

func TestIsPinned_NilMeta(t *testing.T) {
	if IsPinned(nil) {
		t.Error("expected IsPinned to return false for nil map")
	}
}
