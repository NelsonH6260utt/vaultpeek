package vault

import (
	"testing"
)

func TestRollbackResult_Fields(t *testing.T) {
	r := RollbackResult{
		Path:        "secret/myapp/config",
		FromVersion: 5,
		ToVersion:   3,
		Success:     true,
	}

	if r.Path != "secret/myapp/config" {
		t.Errorf("expected path %q, got %q", "secret/myapp/config", r.Path)
	}
	if r.FromVersion != 5 {
		t.Errorf("expected FromVersion 5, got %d", r.FromVersion)
	}
	if r.ToVersion != 3 {
		t.Errorf("expected ToVersion 3, got %d", r.ToVersion)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestRollbackResult_FailureState(t *testing.T) {
	err := fmt.Errorf("vault unreachable")
	r := RollbackResult{
		Path:    "secret/myapp/config",
		Success: false,
		Error:   err,
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Error == nil {
		t.Error("expected non-nil error")
	}
}

func TestRollbackResult_ZeroVersions(t *testing.T) {
	r := RollbackResult{
		Path:        "secret/test",
		FromVersion: 0,
		ToVersion:   0,
		Success:     false,
	}

	if r.FromVersion != 0 || r.ToVersion != 0 {
		t.Errorf("expected zero versions, got from=%d to=%d", r.FromVersion, r.ToVersion)
	}
}

func TestRollbackResult_VersionOrdering(t *testing.T) {
	r := RollbackResult{
		FromVersion: 10,
		ToVersion:   2,
		Success:     true,
	}

	if r.ToVersion >= r.FromVersion {
		t.Errorf("expected ToVersion < FromVersion for a rollback, got %d >= %d", r.ToVersion, r.FromVersion)
	}
}
