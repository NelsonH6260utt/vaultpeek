package vault

import (
	"testing"
)

func TestLockResult_Fields(t *testing.T) {
	r := LockResult{
		Path:    "secret/myapp",
		Locked:  true,
		Success: true,
	}
	if r.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", r.Path)
	}
	if !r.Locked {
		t.Error("expected Locked to be true")
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestLockResult_FailureState(t *testing.T) {
	r := LockResult{
		Path:    "secret/myapp",
		Locked:  false,
		Success: false,
	}
	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Locked {
		t.Error("expected Locked to be false")
	}
}

func TestIsLocked_True(t *testing.T) {
	meta := map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"locked": "true",
		},
	}
	if !IsLocked(meta) {
		t.Error("expected IsLocked to return true")
	}
}

func TestIsLocked_False(t *testing.T) {
	meta := map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"locked": "false",
		},
	}
	if IsLocked(meta) {
		t.Error("expected IsLocked to return false")
	}
}

func TestIsLocked_Missing(t *testing.T) {
	if IsLocked(map[string]interface{}{}) {
		t.Error("expected IsLocked to return false for empty meta")
	}
}

func TestLockResult_NilClient(t *testing.T) {
	r := LockSecret(nil, "secret", "myapp", "admin", "maintenance")
	if r.Success {
		t.Error("expected failure with nil client")
	}
	if r.Error == nil {
		t.Error("expected non-nil error")
	}
}
