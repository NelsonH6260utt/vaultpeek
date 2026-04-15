package output

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/vaultpeek/internal/vault"
)

func TestPrintProtectResult_Success(t *testing.T) {
	var out, errOut bytes.Buffer
	result := vault.ProtectResult{
		Path:    "myapp/config",
		Mount:   "secret",
		Success: true,
	}

	printProtectResultTo(&out, &errOut, result)

	if errOut.Len() != 0 {
		t.Errorf("expected no stderr output, got: %s", errOut.String())
	}
	if !strings.Contains(out.String(), "protected") {
		t.Errorf("expected 'protected' in output, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "secret/myapp/config") {
		t.Errorf("expected path in output, got: %s", out.String())
	}
}

func TestPrintProtectResult_Failure(t *testing.T) {
	var out, errOut bytes.Buffer
	result := vault.ProtectResult{
		Path:    "myapp/config",
		Mount:   "secret",
		Success: false,
		Error:   errors.New("permission denied"),
	}

	printProtectResultTo(&out, &errOut, result)

	if out.Len() != 0 {
		t.Errorf("expected no stdout output on failure, got: %s", out.String())
	}
	if !strings.Contains(errOut.String(), "error") {
		t.Errorf("expected 'error' in stderr output, got: %s", errOut.String())
	}
	if !strings.Contains(errOut.String(), "permission denied") {
		t.Errorf("expected error reason in stderr, got: %s", errOut.String())
	}
}

func TestPrintProtectResult_FailureNoError(t *testing.T) {
	var out, errOut bytes.Buffer
	result := vault.ProtectResult{
		Path:    "myapp/config",
		Mount:   "secret",
		Success: false,
		Error:   nil,
	}

	printProtectResultTo(&out, &errOut, result)

	if out.Len() != 0 {
		t.Errorf("expected no stdout output on failure, got: %s", out.String())
	}
	if !strings.Contains(errOut.String(), "myapp/config") {
		t.Errorf("expected path in error message, got: %s", errOut.String())
	}
}

func TestPrintProtectResult_MetadataHint(t *testing.T) {
	var out, errOut bytes.Buffer
	result := vault.ProtectResult{
		Path:    "db/prod",
		Mount:   "kv",
		Success: true,
	}

	printProtectResultTo(&out, &errOut, result)

	if !strings.Contains(out.String(), "vaultpeek_protected") {
		t.Errorf("expected metadata hint in output, got: %s", out.String())
	}
}
