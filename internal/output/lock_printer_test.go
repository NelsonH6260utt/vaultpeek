package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/vaultpeek/internal/vault"
)

func TestPrintLockResult_Locked(t *testing.T) {
	var buf bytes.Buffer
	r := vault.LockResult{Path: "secret/app", Locked: true, Success: true}
	printLockResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "locked") {
		t.Errorf("expected 'locked' in output, got: %s", out)
	}
	if !strings.Contains(out, "secret/app") {
		t.Errorf("expected path in output, got: %s", out)
	}
}

func TestPrintLockResult_Unlocked(t *testing.T) {
	var buf bytes.Buffer
	r := vault.LockResult{Path: "secret/app", Locked: false, Success: true}
	printLockResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "unlocked") {
		t.Errorf("expected 'unlocked' in output, got: %s", out)
	}
}

func TestPrintLockResult_Failure(t *testing.T) {
	var buf bytes.Buffer
	r := vault.LockResult{Path: "secret/app", Success: false, Error: fmt.Errorf("permission denied")}
	printLockResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "error") {
		t.Errorf("expected 'error' in output, got: %s", out)
	}
	if !strings.Contains(out, "permission denied") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestPrintLockResult_FailureNoError(t *testing.T) {
	var buf bytes.Buffer
	r := vault.LockResult{Path: "secret/x", Success: false}
	printLockResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "error") {
		t.Errorf("expected 'error' in output, got: %s", out)
	}
}
