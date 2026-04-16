package vault

import (
	"errors"
	"testing"
)

func TestCloneResult_Fields(t *testing.T) {
	r := CloneResult{
		SourcePath: "secret/src",
		DestPath:   "secret/dst",
		Keys:       []string{"foo", "bar"},
		Success:    true,
	}
	if r.SourcePath != "secret/src" {
		t.Errorf("unexpected SourcePath: %s", r.SourcePath)
	}
	if r.DestPath != "secret/dst" {
		t.Errorf("unexpected DestPath: %s", r.DestPath)
	}
	if len(r.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Keys))
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestCloneResult_FailureState(t *testing.T) {
	r := CloneResult{
		Success: false,
		Error:   errors.New("clone failed"),
	}
	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Error == nil {
		t.Error("expected non-nil Error")
	}
}

func TestCloneResult_EmptyKeys(t *testing.T) {
	r := CloneResult{Keys: []string{}}
	if len(r.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(r.Keys))
	}
}

func TestCloneResult_NilKeys(t *testing.T) {
	r := CloneResult{}
	if r.Keys != nil {
		t.Error("expected nil Keys")
	}
}

func TestCloneResult_PathsPreserved(t *testing.T) {
	r := CloneResult{SourcePath: "a/b", DestPath: "c/d"}
	if r.SourcePath != "a/b" || r.DestPath != "c/d" {
		t.Error("paths not preserved")
	}
}
