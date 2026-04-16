package output

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/yourusername/vaultpeek/internal/vault"
)

func TestPrintTagResult_Success(t *testing.T) {
	var buf bytes.Buffer
	r := vault.TagResult{
		Path:    "secret/app",
		Tags:    map[string]string{"env": "prod"},
		Success: true,
	}
	printTagResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "Tagged: secret/app") {
		t.Errorf("expected path in output, got: %s", out)
	}
	if !strings.Contains(out, "env = prod") {
		t.Errorf("expected tag in output, got: %s", out)
	}
}

func TestPrintTagResult_Failure(t *testing.T) {
	var buf bytes.Buffer
	r := vault.TagResult{
		Path:    "secret/app",
		Success: false,
		Err:     fmt.Errorf("permission denied"),
	}
	printTagResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "[error]") {
		t.Errorf("expected error marker, got: %s", out)
	}
	if !strings.Contains(out, "permission denied") {
		t.Errorf("expected error message, got: %s", out)
	}
}

func TestPrintTagResult_EmptyTags(t *testing.T) {
	var buf bytes.Buffer
	r := vault.TagResult{
		Path:    "secret/app",
		Tags:    map[string]string{},
		Success: true,
	}
	printTagResultTo(&buf, r)
	out := buf.String()
	if !strings.Contains(out, "(no tags)") {
		t.Errorf("expected no-tags message, got: %s", out)
	}
}

func TestPrintTagResult_SortedKeys(t *testing.T) {
	var buf bytes.Buffer
	r := vault.TagResult{
		Path:    "secret/app",
		Tags:    map[string]string{"zzz": "last", "aaa": "first"},
		Success: true,
	}
	printTagResultTo(&buf, r)
	out := buf.String()
	ia := strings.Index(out, "aaa")
	iz := strings.Index(out, "zzz")
	if ia > iz {
		t.Errorf("expected aaa before zzz in output")
	}
}
