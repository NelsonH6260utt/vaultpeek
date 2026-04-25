package vault

import (
	"errors"
	"testing"
	"time"
)

func TestTouchResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	r := TouchResult{
		Path:      "secret/myapp/config",
		Success:   true,
		Error:     nil,
		TouchedAt: now,
		Version:   5,
	}

	if r.Path != "secret/myapp/config" {
		t.Errorf("expected path 'secret/myapp/config', got %q", r.Path)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if r.Error != nil {
		t.Errorf("expected no error, got %v", r.Error)
	}
	if r.Version != 5 {
		t.Errorf("expected version 5, got %d", r.Version)
	}
	if r.TouchedAt != now {
		t.Errorf("expected TouchedAt %v, got %v", now, r.TouchedAt)
	}
}

func TestTouchResult_FailureState(t *testing.T) {
	expectedErr := errors.New("vault unavailable")
	r := TouchResult{
		Path:    "secret/broken",
		Success: false,
		Error:   expectedErr,
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Error == nil {
		t.Error("expected an error, got nil")
	}
	if r.Error.Error() != "vault unavailable" {
		t.Errorf("unexpected error message: %s", r.Error.Error())
	}
}

func TestTouchResult_ZeroVersion(t *testing.T) {
	r := TouchResult{
		Path:    "secret/new",
		Success: true,
		Version: 0,
	}

	if r.Version != 0 {
		t.Errorf("expected version 0, got %d", r.Version)
	}
}

func TestTouchResult_TouchedAtIsUTC(t *testing.T) {
	now := time.Now().UTC()
	r := TouchResult{
		Path:      "secret/app/key",
		Success:   true,
		TouchedAt: now,
	}

	if r.TouchedAt.Location() != time.UTC {
		t.Errorf("expected UTC location, got %v", r.TouchedAt.Location())
	}
}

func TestTouchResult_PathPreserved(t *testing.T) {
	path := "kv/prod/database/credentials"
	r := TouchResult{
		Path:    path,
		Success: true,
	}

	if r.Path != path {
		t.Errorf("path not preserved: got %q, want %q", r.Path, path)
	}
}
