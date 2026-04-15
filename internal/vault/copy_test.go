package vault

import (
	"testing"
)

func TestCopyResult_Fields(t *testing.T) {
	r := CopyResult{
		SourcePath: "secret/src",
		DestPath:   "secret/dst",
		Keys:       []string{"foo", "bar"},
		Success:    true,
		Err:        nil,
	}

	if r.SourcePath != "secret/src" {
		t.Errorf("expected SourcePath 'secret/src', got %q", r.SourcePath)
	}
	if r.DestPath != "secret/dst" {
		t.Errorf("expected DestPath 'secret/dst', got %q", r.DestPath)
	}
	if len(r.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Keys))
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if r.Err != nil {
		t.Errorf("expected nil Err, got %v", r.Err)
	}
}

func TestCopyResult_FailureState(t *testing.T) {
	r := CopyResult{
		SourcePath: "secret/missing",
		DestPath:   "secret/dst",
		Success:    false,
		Err:        fmt.Errorf("source secret %q not found", "secret/missing"),
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Err == nil {
		t.Error("expected non-nil Err")
	}
}

func TestCopyResult_EmptyKeys(t *testing.T) {
	r := CopyResult{
		Keys: []string{},
	}
	if len(r.Keys) != 0 {
		t.Errorf("expected empty keys slice, got %d elements", len(r.Keys))
	}
}

func TestCopyResult_NilKeys(t *testing.T) {
	r := CopyResult{}
	if r.Keys != nil {
		t.Errorf("expected nil Keys, got %v", r.Keys)
	}
}
