package vault

import (
	"testing"
	"time"
)

func TestRestoreResult_Fields(t *testing.T) {
	now := time.Now()
	r := RestoreResult{
		Path:        "secret/myapp",
		FromVersion: 2,
		NewVersion:  5,
		RestoredAt:  now,
		Keys:        []string{"api_key", "db_pass"},
		Success:     true,
	}

	if r.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", r.Path)
	}
	if r.FromVersion != 2 {
		t.Errorf("expected from version 2, got %d", r.FromVersion)
	}
	if r.NewVersion != 5 {
		t.Errorf("expected new version 5, got %d", r.NewVersion)
	}
	if !r.Success {
		t.Error("expected success to be true")
	}
	if len(r.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Keys))
	}
}

func TestRestoreResult_FailureState(t *testing.T) {
	r := RestoreResult{
		Path:    "secret/missing",
		Success: false,
		Err:     fmt.Errorf("version not found"),
	}

	if r.Success {
		t.Error("expected success to be false")
	}
	if r.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestRestoreResult_ZeroVersion(t *testing.T) {
	r := RestoreResult{
		Path:        "secret/app",
		FromVersion: 0,
		NewVersion:  0,
		Success:     false,
	}

	if r.FromVersion != 0 {
		t.Errorf("expected 0, got %d", r.FromVersion)
	}
	if r.NewVersion != 0 {
		t.Errorf("expected 0, got %d", r.NewVersion)
	}
}

func TestRestoreResult_EmptyKeys(t *testing.T) {
	r := RestoreResult{
		Path:    "secret/empty",
		Keys:    []string{},
		Success: true,
	}

	if len(r.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(r.Keys))
	}
}
