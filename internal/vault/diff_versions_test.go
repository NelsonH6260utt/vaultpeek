package vault

import (
	"testing"
)

func TestVersionDiffResult_Fields(t *testing.T) {
	r := VersionDiffResult{
		Path:     "secret/myapp",
		VersionA: 1,
		VersionB: 2,
		Added:    map[string]string{"newkey": "newval"},
		Removed:  map[string]string{"oldkey": "oldval"},
		Changed:  map[string][2]string{"key": {"before", "after"}},
		Success:  true,
	}
	if r.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", r.Path)
	}
	if r.VersionA != 1 || r.VersionB != 2 {
		t.Errorf("unexpected versions: %d %d", r.VersionA, r.VersionB)
	}
	if len(r.Added) != 1 {
		t.Errorf("expected 1 added key")
	}
	if len(r.Removed) != 1 {
		t.Errorf("expected 1 removed key")
	}
	if len(r.Changed) != 1 {
		t.Errorf("expected 1 changed key")
	}
}

func TestVersionDiffResult_FailureState(t *testing.T) {
	r := VersionDiffResult{
		Success: false,
		Err:     fmt.Errorf("vault unreachable"),
	}
	if r.Success {
		t.Error("expected failure state")
	}
	if r.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestVersionDiffResult_EmptyMaps(t *testing.T) {
	r := VersionDiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
		Success: true,
	}
	if len(r.Added) != 0 || len(r.Removed) != 0 || len(r.Changed) != 0 {
		t.Error("expected all maps to be empty")
	}
}

func TestVersionDiffResult_ChangedValues(t *testing.T) {
	r := VersionDiffResult{
		Changed: map[string][2]string{
			"password": {"hunter2", "correct-horse"},
		},
		Success: true,
	}
	pair, ok := r.Changed["password"]
	if !ok {
		t.Fatal("expected changed key 'password'")
	}
	if pair[0] != "hunter2" || pair[1] != "correct-horse" {
		t.Errorf("unexpected changed values: %v", pair)
	}
}
