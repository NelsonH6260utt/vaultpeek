package vault

import (
	"testing"
)

func TestTemplateResult_Fields(t *testing.T) {
	r := TemplateResult{
		Path:    "secret/app/config",
		Success: true,
		Rendered: map[string]string{"host": "localhost"},
		Missing:  []string{},
	}
	if r.Path != "secret/app/config" {
		t.Errorf("expected path, got %q", r.Path)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestTemplateResult_FailureState(t *testing.T) {
	r := TemplateResult{Success: false, Err: fmt.Errorf("boom")}
	if r.Success {
		t.Error("expected failure state")
	}
}

func TestRenderTemplate_AllResolved(t *testing.T) {
	data := map[string]interface{}{"host": "db.internal", "port": "5432"}
	res := RenderTemplate("secret/app", "connect to {{host}}:{{port}}", data)
	if !res.Success {
		t.Fatalf("expected success, missing: %v", res.Missing)
	}
	if len(res.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", res.Missing)
	}
	if res.Rendered["host"] != "db.internal" {
		t.Errorf("unexpected rendered host: %s", res.Rendered["host"])
	}
}

func TestRenderTemplate_MissingKey(t *testing.T) {
	data := map[string]interface{}{"host": "db.internal"}
	res := RenderTemplate("secret/app", "{{host}}:{{port}}", data)
	if res.Success {
		t.Error("expected failure due to missing key")
	}
	if len(res.Missing) != 1 || res.Missing[0] != "port" {
		t.Errorf("expected missing key 'port', got %v", res.Missing)
	}
}

func TestRenderTemplate_NilData(t *testing.T) {
	res := RenderTemplate("secret/app", "{{key}}", nil)
	if res.Success {
		t.Error("expected failure for nil data")
	}
	if res.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestRenderTemplate_NoPlaceholders(t *testing.T) {
	data := map[string]interface{}{"x": "y"}
	res := RenderTemplate("secret/app", "static string", data)
	if !res.Success {
		t.Errorf("expected success for static string, got missing: %v", res.Missing)
	}
}

func init() {
	// pull in fmt for TestTemplateResult_FailureState
	_ = fmt.Sprintf
}
