package vault

import (
	"testing"
	"time"
)

func TestRotateResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	r := RotateResult{
		Path:       "secret/myapp/db",
		OldVersion: 3,
		NewVersion: 4,
		RotatedAt:  now,
		Keys:       []string{"password", "username"},
		Success:    true,
	}

	if r.Path != "secret/myapp/db" {
		t.Errorf("expected path secret/myapp/db, got %s", r.Path)
	}
	if r.OldVersion != 3 {
		t.Errorf("expected old version 3, got %d", r.OldVersion)
	}
	if r.NewVersion != 4 {
		t.Errorf("expected new version 4, got %d", r.NewVersion)
	}
	if !r.Success {
		t.Error("expected success to be true")
	}
	if len(r.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Keys))
	}
}

func TestRotateResult_FailureState(t *testing.T) {
	r := RotateResult{
		Path:    "secret/missing",
		Success: false,
		Err:     fmt.Errorf("rotate: no data at secret/missing"),
	}

	if r.Success {
		t.Error("expected success to be false")
	}
	if r.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestRotateResult_ZeroVersions(t *testing.T) {
	r := RotateResult{}
	if r.OldVersion != 0 {
		t.Errorf("expected zero old version, got %d", r.OldVersion)
	}
	if r.NewVersion != 0 {
		t.Errorf("expected zero new version, got %d", r.NewVersion)
	}
}

func TestRotateOptions_NilTransform(t *testing.T) {
	opts := RotateOptions{}
	if opts.Transform != nil {
		t.Error("expected nil transform by default")
	}
}

func TestRotateOptions_TransformApplied(t *testing.T) {
	called := false
	opts := RotateOptions{
		Transform: func(key, value string) string {
			called = true
			return "rotated:" + value
		},
	}

	result := opts.Transform("password", "secret123")
	if !called {
		t.Error("expected transform to be called")
	}
	if result != "rotated:secret123" {
		t.Errorf("unexpected transform result: %s", result)
	}
}
