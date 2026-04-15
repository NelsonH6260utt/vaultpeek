package vault

import (
	"testing"
)

func TestRenameResult_Fields(t *testing.T) {
	r := RenameResult{
		SourcePath: "old/path",
		DestPath:   "new/path",
		Keys:       []string{"username", "password"},
		Success:    true,
	}

	if r.SourcePath != "old/path" {
		t.Errorf("expected SourcePath 'old/path', got %q", r.SourcePath)
	}
	if r.DestPath != "new/path" {
		t.Errorf("expected DestPath 'new/path', got %q", r.DestPath)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if len(r.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Keys))
	}
}

func TestRenameResult_FailureState(t *testing.T) {
	r := RenameResult{
		SourcePath: "src",
		DestPath:   "dst",
		Success:    false,
		Err:        fmt.Errorf("source secret not found"),
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Err == nil {
		t.Error("expected non-nil Err")
	}
}

func TestRenameResult_EmptyKeys(t *testing.T) {
	r := RenameResult{
		SourcePath: "a",
		DestPath:   "b",
		Keys:       []string{},
		Success:    true,
	}

	if len(r.Keys) != 0 {
		t.Errorf("expected 0 keys, got %d", len(r.Keys))
	}
}

func TestRenameResult_NilKeys(t *testing.T) {
	r := RenameResult{
		SourcePath: "a",
		DestPath:   "b",
	}

	if r.Keys != nil {
		t.Errorf("expected nil Keys, got %v", r.Keys)
	}
	if r.Success {
		t.Error("expected Success to default to false")
	}
}
