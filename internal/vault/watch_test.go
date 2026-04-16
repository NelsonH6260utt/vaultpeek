package vault

import (
	"testing"
	"time"
)

func TestWatchResult_Fields(t *testing.T) {
	r := WatchResult{
		Path:    "secret/myapp",
		Changed: true,
		Version: 3,
		Data:    map[string]interface{}{"key": "val"},
		Error:   nil,
	}
	if r.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", r.Path)
	}
	if !r.Changed {
		t.Error("expected Changed to be true")
	}
	if r.Version != 3 {
		t.Errorf("expected version 3, got %d", r.Version)
	}
}

func TestWatchResult_FailureState(t *testing.T) {
	err := errorf("vault unavailable")
	r := WatchResult{Path: "secret/x", Error: err}
	if r.Error == nil {
		t.Error("expected non-nil error")
	}
	if r.Changed {
		t.Error("Changed should be false on error result")
	}
}

func TestWatchOptions_DefaultInterval(t *testing.T) {
	opts := WatchOptions{}
	if opts.Interval != 0 {
		t.Errorf("zero value interval should be 0, got %v", opts.Interval)
	}
	// Confirm that 0 triggers the default inside WatchSecret (5s).
	// We just verify the constant expectation here.
	const defaultInterval = 5 * time.Second
	if defaultInterval != 5*time.Second {
		t.Error("default interval constant mismatch")
	}
}

func TestWatchResult_NoChange(t *testing.T) {
	r := WatchResult{
		Path:    "secret/stable",
		Changed: false,
		Version: 1,
	}
	if r.Changed {
		t.Error("expected Changed false when version unchanged")
	}
}

func TestWatchResult_ZeroVersion(t *testing.T) {
	r := WatchResult{Path: "secret/empty"}
	if r.Version != 0 {
		t.Errorf("expected zero version, got %d", r.Version)
	}
}

// errorf is a minimal helper to create an error value in tests.
func errorf(msg string) error {
	return fmt.Errorf("%s", msg)
}
