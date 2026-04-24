package vault

import (
	"testing"
)

func TestLabelResult_Fields(t *testing.T) {
	r := LabelResult{
		Path:    "secret/data/myapp/config",
		Labels:  map[string]string{"env": "prod", "team": "platform"},
		Success: true,
	}

	if r.Path != "secret/data/myapp/config" {
		t.Errorf("expected path, got %q", r.Path)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if r.Labels["env"] != "prod" {
		t.Errorf("expected label env=prod, got %q", r.Labels["env"])
	}
}

func TestLabelResult_FailureState(t *testing.T) {
	r := LabelResult{
		Path:    "secret/data/myapp/config",
		Success: false,
		Err:     fmt.Errorf("permission denied"),
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Err == nil {
		t.Error("expected non-nil Err")
	}
}

func TestSortedLabelKeys_Order(t *testing.T) {
	labels := map[string]string{
		"team":    "platform",
		"env":     "staging",
		"version": "v2",
	}
	keys := SortedLabelKeys(labels)
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "env" || keys[1] != "team" || keys[2] != "version" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSortedLabelKeys_Empty(t *testing.T) {
	keys := SortedLabelKeys(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}

func TestLabelResult_NilLabels(t *testing.T) {
	r := LabelResult{
		Path:    "secret/data/app",
		Labels:  nil,
		Success: false,
	}

	if r.Labels != nil {
		t.Error("expected nil Labels")
	}
	keys := SortedLabelKeys(map[string]string{})
	if len(keys) != 0 {
		t.Error("SortedLabelKeys should handle empty map")
	}
}
