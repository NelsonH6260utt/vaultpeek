package vault

import (
	"errors"
	"testing"
)

func TestPromoteResult_Fields(t *testing.T) {
	r := PromoteResult{
		SourcePath: "staging/db",
		DestPath:   "production/db",
		Keys:       []string{"password", "user"},
		Success:    true,
	}

	if r.SourcePath != "staging/db" {
		t.Errorf("expected SourcePath staging/db, got %s", r.SourcePath)
	}
	if r.DestPath != "production/db" {
		t.Errorf("expected DestPath production/db, got %s", r.DestPath)
	}
	if len(r.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Keys))
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestPromoteResult_FailureState(t *testing.T) {
	r := PromoteResult{
		SourcePath: "staging/db",
		DestPath:   "production/db",
		Success:    false,
		Err:        errors.New("permission denied"),
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Err == nil {
		t.Error("expected Err to be set")
	}
}

func TestPromoteResult_EmptyKeys(t *testing.T) {
	r := PromoteResult{
		SourcePath: "staging/app",
		DestPath:   "production/app",
		Keys:       []string{},
		Success:    true,
	}

	if len(r.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(r.Keys))
	}
}

func TestPromoteResult_NilKeys(t *testing.T) {
	r := PromoteResult{
		SourcePath: "staging/app",
		DestPath:   "production/app",
		Success:    false,
	}

	if r.Keys != nil {
		t.Errorf("expected nil Keys, got %v", r.Keys)
	}
}

func TestPromoteResult_PathsPreserved(t *testing.T) {
	src := "env/staging/service/config"
	dest := "env/production/service/config"

	r := PromoteResult{SourcePath: src, DestPath: dest}

	if r.SourcePath != src {
		t.Errorf("source path mismatch: got %s", r.SourcePath)
	}
	if r.DestPath != dest {
		t.Errorf("dest path mismatch: got %s", r.DestPath)
	}
}
