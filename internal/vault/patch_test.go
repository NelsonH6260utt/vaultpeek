package vault

import (
	"testing"
)

func TestPatchResult_Fields(t *testing.T) {
	r := PatchResult{
		Path:        "secret/myapp/config",
		UpdatedKeys: []string{"db_pass", "api_key"},
		Success:     true,
	}

	if r.Path != "secret/myapp/config" {
		t.Errorf("expected path 'secret/myapp/config', got %q", r.Path)
	}
	if len(r.UpdatedKeys) != 2 {
		t.Errorf("expected 2 updated keys, got %d", len(r.UpdatedKeys))
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestPatchResult_FailureState(t *testing.T) {
	r := PatchResult{
		Path:    "secret/myapp/config",
		Success: false,
		Error:   fmt.Errorf("connection refused"),
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Error == nil {
		t.Error("expected non-nil Error")
	}
}

func TestPatchResult_EmptyKeys(t *testing.T) {
	r := PatchResult{
		Path:        "secret/empty",
		UpdatedKeys: []string{},
		Success:     true,
	}

	if len(r.UpdatedKeys) != 0 {
		t.Errorf("expected 0 updated keys, got %d", len(r.UpdatedKeys))
	}
}

func TestPatchResult_NilKeys(t *testing.T) {
	r := PatchResult{
		Path:    "secret/nil-keys",
		Success: true,
	}

	if r.UpdatedKeys != nil {
		t.Errorf("expected nil UpdatedKeys, got %v", r.UpdatedKeys)
	}
}
