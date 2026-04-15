package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your-org/vaultpeek/internal/vault"
)

func TestPrintRenameResult_Success(t *testing.T) {
	result := vault.RenameResult{
		Success:    true,
		SourcePath: "secret/old",
		DestPath:   "secret/new",
		Keys:       []string{"username", "password"},
	}

	var buf bytes.Buffer
	printRenameResultTo(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "✓") {
		t.Error("expected success symbol in output")
	}
	if !strings.Contains(out, "secret/old") {
		t.Error("expected source path in output")
	}
	if !strings.Contains(out, "secret/new") {
		t.Error("expected dest path in output")
	}
	if !strings.Contains(out, "username") {
		t.Error("expected key 'username' in output")
	}
}

func TestPrintRenameResult_Failure(t *testing.T) {
	result := vault.RenameResult{
		Success: false,
		Error:   "permission denied",
	}

	var buf bytes.Buffer
	printRenameResultTo(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "✗") {
		t.Error("expected failure symbol in output")
	}
	if !strings.Contains(out, "permission denied") {
		t.Error("expected error message in output")
	}
}

func TestPrintRenameResult_EmptyKeys(t *testing.T) {
	result := vault.RenameResult{
		Success:    true,
		SourcePath: "secret/a",
		DestPath:   "secret/b",
		Keys:       []string{},
	}

	var buf bytes.Buffer
	printRenameResultTo(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "(none)") {
		t.Error("expected '(none)' when no keys are present")
	}
}

func TestPrintRenameResult_KeyCount(t *testing.T) {
	result := vault.RenameResult{
		Success:    true,
		SourcePath: "secret/x",
		DestPath:   "secret/y",
		Keys:       []string{"a", "b", "c"},
	}

	var buf bytes.Buffer
	printRenameResultTo(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "3 transferred") {
		t.Errorf("expected key count in output, got: %s", out)
	}
}
