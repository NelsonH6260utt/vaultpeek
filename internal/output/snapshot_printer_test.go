package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourusername/vaultpeek/internal/vault"
)

func makeSnapshot(success bool) vault.SnapshotResult {
	return vault.SnapshotResult{
		TakenAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Success: success,
		Entries: []vault.SnapshotEntry{
			{Path: "secret/app", Data: map[string]interface{}{"db": "postgres", "port": "5432"}},
		},
		Err: nil,
	}
}

func TestPrintSnapshotResult_Success(t *testing.T) {
	var buf bytes.Buffer
	printSnapshotResultTo(&buf, makeSnapshot(true))
	out := buf.String()
	if !strings.Contains(out, "secret/app") {
		t.Error("expected path in output")
	}
	if !strings.Contains(out, "Paths captured: 1") {
		t.Error("expected count in output")
	}
}

func TestPrintSnapshotResult_Failure(t *testing.T) {
	var buf bytes.Buffer
	r := vault.SnapshotResult{Success: false}
	printSnapshotResultTo(&buf, r)
	if !strings.Contains(buf.String(), "[error]") {
		t.Error("expected error marker in output")
	}
}

func TestPrintSnapshotResult_SortedKeys(t *testing.T) {
	var buf bytes.Buffer
	printSnapshotResultTo(&buf, makeSnapshot(true))
	out := buf.String()
	dbIdx := strings.Index(out, "db")
	portIdx := strings.Index(out, "port")
	if dbIdx > portIdx {
		t.Error("expected keys to be sorted alphabetically")
	}
}

func TestPrintSnapshotResult_ShowsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	printSnapshotResultTo(&buf, makeSnapshot(true))
	if !strings.Contains(buf.String(), "2024-06-01") {
		t.Error("expected timestamp in output")
	}
}
