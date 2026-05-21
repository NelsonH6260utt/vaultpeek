package vault

import (
	"testing"
)

func TestAliasResult_Fields(t *testing.T) {
	r := AliasResult{Name: "prod", Path: "secret/prod/app", Success: true}
	if r.Name != "prod" {
		t.Errorf("expected Name=prod, got %s", r.Name)
	}
	if r.Path != "secret/prod/app" {
		t.Errorf("expected Path=secret/prod/app, got %s", r.Path)
	}
	if !r.Success {
		t.Error("expected Success=true")
	}
}

func TestAliasResult_FailureState(t *testing.T) {
	r := AliasResult{Success: false, Err: errorf("something went wrong")}
	if r.Success {
		t.Error("expected Success=false")
	}
	if r.Err == nil {
		t.Error("expected non-nil Err")
	}
}

func TestAliasStore_SetAndResolve(t *testing.T) {
	s := NewAliasStore()
	res := s.Set("dev", "secret/dev/config")
	if !res.Success {
		t.Fatalf("expected success, got err: %v", res.Err)
	}
	path, ok := s.Resolve("dev")
	if !ok {
		t.Fatal("expected alias to resolve")
	}
	if path != "secret/dev/config" {
		t.Errorf("expected path=secret/dev/config, got %s", path)
	}
}

func TestAliasStore_Remove(t *testing.T) {
	s := NewAliasStore()
	s.Set("staging", "secret/staging/app")
	res := s.Remove("staging")
	if !res.Success {
		t.Fatalf("expected success, got err: %v", res.Err)
	}
	if _, ok := s.Resolve("staging"); ok {
		t.Error("expected alias to be removed")
	}
}

func TestAliasStore_RemoveMissing(t *testing.T) {
	s := NewAliasStore()
	res := s.Remove("nonexistent")
	if res.Success {
		t.Error("expected failure for missing alias")
	}
	if res.Err == nil {
		t.Error("expected non-nil Err")
	}
}

func TestAliasStore_ListSorted(t *testing.T) {
	s := NewAliasStore()
	s.Set("zebra", "secret/z")
	s.Set("alpha", "secret/a")
	s.Set("mango", "secret/m")
	list := s.List()
	if len(list) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(list))
	}
	if list[0].Name != "alpha" || list[1].Name != "mango" || list[2].Name != "zebra" {
		t.Errorf("unexpected sort order: %v", list)
	}
}

func TestAliasStore_SetEmptyName(t *testing.T) {
	s := NewAliasStore()
	res := s.Set("", "secret/path")
	if res.Success {
		t.Error("expected failure for empty name")
	}
}

func errorf(msg string) error {
	return fmt.Errorf("%s", msg)
}
