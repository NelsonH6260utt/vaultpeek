package vault

import (
	"testing"
)

func TestLintResult_Fields(t *testing.T) {
	r := LintResult{Path: "secret/foo", Success: true, Warnings: []string{"w1"}}
	if r.Path != "secret/foo" {
		t.Errorf("expected path secret/foo, got %s", r.Path)
	}
	if !r.Success {
		t.Error("expected success true")
	}
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}

func TestLintSecret_NilData(t *testing.T) {
	r := LintSecret("secret/foo", nil, DefaultLintRules())
	if r.Success {
		t.Error("expected failure for nil data")
	}
	if r.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestLintSecret_EmptyValue(t *testing.T) {
	data := map[string]interface{}{"username": "", "password": "s3cr3t"}
	r := LintSecret("secret/foo", data, DefaultLintRules())
	if !r.Success {
		t.Error("expected success")
	}
	if len(r.Warnings) == 0 {
		t.Error("expected at least one warning for empty value")
	}
}

func TestLintSecret_NoWarnings(t *testing.T) {
	data := map[string]interface{}{"username": "admin", "password": "s3cr3t"}
	r := LintSecret("secret/foo", data, DefaultLintRules())
	if !r.Success {
		t.Error("expected success")
	}
	if len(r.Warnings) != 0 {
		t.Errorf("expected no warnings, got %d", len(r.Warnings))
	}
}

func TestLintSecret_CustomRule(t *testing.T) {
	rule := func(key, value string) string {
		if key == "forbidden" {
			return "key 'forbidden' is not allowed"
		}
		return ""
	}
	data := map[string]interface{}{"forbidden": "value"}
	r := LintSecret("secret/bar", data, []LintRule{rule})
	if len(r.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(r.Warnings))
	}
}
