package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/robbrockway/vaultpeek/internal/vault"
)

func TestPrintExportResult_Success(t *testing.T) {
	var buf bytes.Buffer
	r := vault.ExportResult{
		Path:    "secret/app/config",
		Data:    map[string]interface{}{"db_pass": "s3cr3t", "api_key": "abc123"},
		Success: true,
	}
	printExportResultTo(&buf, r, false)
	out := buf.String()
	if !strings.Contains(out, "secret/app/config") {
		t.Error("expected path in output")
	}
	if !strings.Contains(out, "db_pass") {
		t.Error("expected key db_pass in output")
	}
	if !strings.Contains(out, "api_key") {
		t.Error("expected key api_key in output")
	}
}

func TestPrintExportResult_Failure(t *testing.T) {
	var buf bytes.Buffer
	r := vault.ExportResult{
		Path:    "bad/path",
		Success: false,
		Err:     fmt.Errorf("not found"),
	}
	printExportResultTo(&buf, r, false)
	out := buf.String()
	if !strings.Contains(out, "error") {
		t.Error("expected error in output")
	}
}

func TestPrintExportResult_JSONMode(t *testing.T) {
	var buf bytes.Buffer
	r := vault.ExportResult{
		Path:    "a/b",
		Data:    map[string]interface{}{"x": "y"},
		Success: true,
	}
	printExportResultTo(&buf, r, true)
	out := buf.String()
	if !strings.Contains(out, "{") {
		t.Error("expected JSON object in output")
	}
	if !strings.Contains(out, `"x"`) {
		t.Error("expected key x in JSON output")
	}
}

func TestPrintExportResult_SortedKeys(t *testing.T) {
	var buf bytes.Buffer
	r := vault.ExportResult{
		Path: "sorted",
		Data: map[string]interface{}{"z_key": 1, "a_key": 2, "m_key": 3},
		Success: true,
	}
	printExportResultTo(&buf, r, false)
	out := buf.String()
	aIdx := strings.Index(out, "a_key")
	mIdx := strings.Index(out, "m_key")
	zIdx := strings.Index(out, "z_key")
	if !(aIdx < mIdx && mIdx < zIdx) {
		t.Error("expected keys in sorted order")
	}
}
