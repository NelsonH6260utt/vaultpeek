package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/vaultpeek/internal/vault"
)

func TestPrintRestoreResult_Success(t *testing.T) {
	r := vault.RestoreResult{
		Path:        "secret/myapp",
		FromVersion: 3,
		NewVersion:  6,
		RestoredAt:  time.Date(2024, 5, 1, 12, 0, 0, 0, time.UTC),
		Keys:        []string{"token", "password"},
		Success:     true,
	}

	var buf bytes.Buffer
	printRestoreResultTo(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "[restored]") {
		t.Error("expected [restored] label")
	}
	if !strings.Contains(out, "secret/myapp") {
		t.Error("expected path in output")
	}
	if !strings.Contains(out, "from version : 3") {
		t.Error("expected from version")
	}
	if !strings.Contains(out, "new version  : 6") {
		t.Error("expected new version")
	}
}

func TestPrintRestoreResult_Failure(t *testing.T) {
	r := vault.RestoreResult{
		Path:    "secret/gone",
		Success: false,
		Err:     fmt.Errorf("version not found"),
	}

	var buf bytes.Buffer
	printRestoreResultTo(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "[error]") {
		t.Error("expected [error] label")
	}
	if !strings.Contains(out, "version not found") {
		t.Error("expected error message")
	}
}

func TestPrintRestoreResult_SortedKeys(t *testing.T) {
	r := vault.RestoreResult{
		Path:        "secret/app",
		FromVersion: 1,
		NewVersion:  2,
		RestoredAt:  time.Now(),
		Keys:        []string{"zebra", "alpha", "mango"},
		Success:     true,
	}

	var buf bytes.Buffer
	printRestoreResultTo(&buf, r)
	out := buf.String()

	alpha := strings.Index(out, "alpha")
	mango := strings.Index(out, "mango")
	zebra := strings.Index(out, "zebra")

	if alpha > mango || mango > zebra {
		t.Error("expected keys to be sorted alphabetically")
	}
}

func TestPrintRestoreResult_EmptyKeys(t *testing.T) {
	r := vault.RestoreResult{
		Path:        "secret/empty",
		FromVersion: 1,
		NewVersion:  2,
		RestoredAt:  time.Now(),
		Keys:        []string{},
		Success:     true,
	}

	var buf bytes.Buffer
	printRestoreResultTo(&buf, r)
	out := buf.String()

	if !strings.Contains(out, "(none)") {
		t.Error("expected (none) for empty keys")
	}
}
