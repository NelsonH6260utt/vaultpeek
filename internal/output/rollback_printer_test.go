package output

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yourusername/vaultpeek/internal/vault"
)

func TestPrintRollbackResult_Success(t *testing.T) {
	var buf bytes.Buffer
	result := vault.RollbackResult{
		Path:        "secret/app/db",
		FromVersion: 4,
		ToVersion:   2,
		Success:     true,
	}

	printRollbackResultTo(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "secret/app/db") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "4") {
		t.Errorf("expected FromVersion 4 in output, got: %s", out)
	}
	if !strings.Contains(out, "2") {
		t.Errorf("expected ToVersion 2 in output, got: %s", out)
	}
	if !strings.Contains(out, "✔") {
		t.Errorf("expected success symbol in output, got: %s", out)
	}
}

func TestPrintRollbackResult_Failure(t *testing.T) {
	var buf bytes.Buffer
	result := vault.RollbackResult{
		Path:    "secret/app/db",
		Success: false,
		Error:   fmt.Errorf("version not found"),
	}

	printRollbackResultTo(&buf, result)
	out := buf.String()

	if !strings.Contains(out, "✗") {
		t.Errorf("expected failure symbol in output, got: %s", out)
	}
	if !strings.Contains(out, "version not found") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestPrintRollbackResult_NewVersionIncrement(t *testing.T) {
	var buf bytes.Buffer
	result := vault.RollbackResult{
		Path:        "secret/svc/token",
		FromVersion: 7,
		ToVersion:   3,
		Success:     true,
	}

	printRollbackResultTo(&buf, result)
	out := buf.String()

	// New version should be FromVersion + 1 = 8
	if !strings.Contains(out, "8") {
		t.Errorf("expected new version 8 in output, got: %s", out)
	}
}
