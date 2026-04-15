package vault

import (
	"errors"
	"testing"
)

func TestProtectResult_Fields(t *testing.T) {
	r := ProtectResult{
		Path:    "secret/myapp",
		Mount:   "secret",
		Success: true,
		Error:   nil,
	}

	if r.Path != "secret/myapp" {
		t.Errorf("expected path %q, got %q", "secret/myapp", r.Path)
	}
	if r.Mount != "secret" {
		t.Errorf("expected mount %q, got %q", "secret", r.Mount)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if r.Error != nil {
		t.Errorf("expected no error, got %v", r.Error)
	}
}

func TestProtectResult_FailureState(t *testing.T) {
	err := errors.New("vault unreachable")
	r := ProtectResult{
		Path:    "secret/myapp",
		Mount:   "secret",
		Success: false,
		Error:   err,
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Error == nil {
		t.Error("expected an error, got nil")
	}
	if !errors.Is(r.Error, err) {
		t.Errorf("unexpected error: %v", r.Error)
	}
}

func TestProtectResult_ZeroValue(t *testing.T) {
	var r ProtectResult

	if r.Success {
		t.Error("zero-value Success should be false")
	}
	if r.Error != nil {
		t.Error("zero-value Error should be nil")
	}
	if r.Path != "" {
		t.Error("zero-value Path should be empty")
	}
}

func TestProtectResult_PathPreserved(t *testing.T) {
	const wantPath = "kv/prod/database"
	r := ProtectResult{Path: wantPath, Mount: "kv", Success: true}
	if r.Path != wantPath {
		t.Errorf("expected path %q, got %q", wantPath, r.Path)
	}
}
