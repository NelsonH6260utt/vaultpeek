package vault

import (
	"testing"
	"time"
)

func TestExpireResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	r := ExpireResult{
		Path:      "secret/myapp",
		TTL:       24 * time.Hour,
		ExpiresAt: now,
		Success:   true,
	}

	if r.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", r.Path)
	}
	if r.TTL != 24*time.Hour {
		t.Errorf("unexpected TTL: %v", r.TTL)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if r.Err != nil {
		t.Errorf("unexpected error: %v", r.Err)
	}
}

func TestExpireResult_FailureState(t *testing.T) {
	r := ExpireResult{
		Path:    "secret/myapp",
		Success: false,
		Err:     fmt.Errorf("something went wrong"),
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestIsExpired_Past(t *testing.T) {
	past := time.Now().UTC().Add(-1 * time.Hour)
	if !IsExpired(past) {
		t.Error("expected past time to be expired")
	}
}

func TestIsExpired_Future(t *testing.T) {
	future := time.Now().UTC().Add(1 * time.Hour)
	if IsExpired(future) {
		t.Error("expected future time to not be expired")
	}
}

func TestIsExpired_ZeroTime(t *testing.T) {
	if IsExpired(time.Time{}) {
		t.Error("expected zero time to not be considered expired")
	}
}

func TestExpireOptions_ZeroTTL(t *testing.T) {
	opts := ExpireOptions{TTL: 0}
	if opts.TTL > 0 {
		t.Error("expected zero TTL")
	}
}
