package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/vaultpeek/internal/vault"
)

func TestPrintCloneResult_Success(t *testing.T) {
	r := vault.CloneResult{
		SourcePath: "secret/src",
		DestPath:   "secret/dst",
		Keys:       []string{"alpha", "beta"},
		Success:    true,
	}
	var buf bytes.Buffer
	printCloneResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "secret/src") {
		t.Error("expected source path in output")
	}
	if !strings.Contains(out, "secret/dst") {
		t.Error("expected dest path in output")
	}
	if !strings.Contains(out, "alpha") {
		t.Error("expected key 'alpha' in output")
	}
}

func TestPrintCloneResult_Failure(t *testing.T) {
	import_err := func(s string) error { return &cloneErr{s} }
	r := vault.CloneResult{
		Success: false,
		Error:   import_err("vault unreachable"),
	}
	var buf bytes.Buffer
	printCloneResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "vault unreachable") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestPrintCloneResult_FailureNoError(t *testing.T) {
	r := vault.CloneResult{Success: false}
	var buf bytes.Buffer
	printCloneResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "unknown error") {
		t.Errorf("expected unknown error message, got: %s", out)
	}
}

func TestPrintCloneResult_SortedKeys(t *testing.T) {
	r := vault.CloneResult{
		SourcePath: "s",
		DestPath:   "d",
		Keys:       []string{"zebra", "apple", "mango"},
		Success:    true,
	}
	var buf bytes.Buffer
	printCloneResultTo(&buf, r)
	out := buf.String()
	appleIdx := strings.Index(out, "apple")
	mangoIdx := strings.Index(out, "mango")
	zebraIdx := strings.Index(out, "zebra")
	if !(appleIdx < mangoIdx && mangoIdx < zebraIdx) {
		t.Error("expected sorted key output")
	}
}

type cloneErr struct{ msg string }

func (e *cloneErr) Error() string { return e.msg }
