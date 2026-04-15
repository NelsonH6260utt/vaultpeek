package vault

import (
	"testing"
)

func TestExportResult_Fields(t *testing.T) {
	r := ExportResult{
		Path:    "secret/myapp/config",
		Data:    map[string]interface{}{"key": "value"},
		Success: true,
	}
	if r.Path != "secret/myapp/config" {
		t.Errorf("expected path, got %s", r.Path)
	}
	if r.Data["key"] != "value" {
		t.Errorf("expected key=value")
	}
}

func TestExportResult_FailureState(t *testing.T) {
	r := ExportResult{Path: "bad/path", Success: false}
	if r.Success {
		t.Error("expected failure state")
	}
	if r.Data != nil {
		t.Error("expected nil data on failure")
	}
}

func TestMarshalExport_Success(t *testing.T) {
	r := ExportResult{
		Path:    "a/b",
		Data:    map[string]interface{}{"foo": "bar"},
		Success: true,
	}
	bytes, err := MarshalExport(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bytes) == 0 {
		t.Error("expected non-empty JSON output")
	}
}

func TestMarshalExport_Failure(t	r := ExportResult{Path: "x", Success: false, Err: errTest}
	_, err := MarshalExport(r)
	if err == nil {
		t.Error("expected error from failed result")
	}
}

func TestMarshalExport_EmptyData(t *testing.T) {
	r := ExportResult{
		Path:    "empty",
		Data:    map[string]interface{}{},
		Success: true,
	}
	bytes, err := MarshalExport(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(bytes) != "{}" {
		t.Errorf("expected {}, got %s", string(bytes))
	}
}

// errTest is a sentinel error for tests.
var errTest = fmt.Errorf("test error")
