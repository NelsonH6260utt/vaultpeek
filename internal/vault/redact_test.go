package vault

import (
	"testing"
)

func TestRedactResult_Fields(t *testing.T) {
	r := RedactResult{
		Path:    "secret/app/config",
		Redacted: map[string]string{"key": "value"},
		Count:   1,
		Success: true,
	}
	if r.Path != "secret/app/config" {
		t.Errorf("expected path, got %q", r.Path)
	}
	if r.Count != 1 {
		t.Errorf("expected count 1, got %d", r.Count)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestRedactResult_FailureState(t *testing.T) {
	r := RedactResult{
		Path:    "secret/app/config",
		Success: false,
		Err:     fmt.Errorf("vault unavailable"),
	}
	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestRedactSecret_NilData(t *testing.T) {
	result := RedactSecret("secret/app", nil, DefaultRedactRules())
	if !result.Success {
		t.Error("expected Success for nil data")
	}
	if len(result.Redacted) != 0 {
		t.Errorf("expected empty redacted map, got %d entries", len(result.Redacted))
	}
}

func TestRedactSecret_MasksPassword(t *testing.T) {
	data := map[string]interface{}{
		"username": "admin",
		"password": "s3cr3t!",
	}
	result := RedactSecret("secret/app", data, DefaultRedactRules())
	if result.Redacted["password"] != "***REDACTED***" {
		t.Errorf("expected password to be redacted, got %q", result.Redacted["password"])
	}
	if result.Redacted["username"] != "admin" {
		t.Errorf("expected username to be unchanged, got %q", result.Redacted["username"])
	}
	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}
}

func TestRedactSecret_NoMatch(t *testing.T) {
	data := map[string]interface{}{
		"region": "us-east-1",
		"env":    "production",
	}
	result := RedactSecret("secret/app", data, DefaultRedactRules())
	if result.Count != 0 {
		t.Errorf("expected count 0, got %d", result.Count)
	}
	if result.Redacted["region"] != "us-east-1" {
		t.Errorf("unexpected value for region: %q", result.Redacted["region"])
	}
}

func TestDefaultRedactRules_Count(t *testing.T) {
	rules := DefaultRedactRules()
	if len(rules) == 0 {
		t.Error("expected at least one default redact rule")
	}
	for _, r := range rules {
		if r.Name == "" {
			t.Error("expected rule to have a non-empty name")
		}
		if r.Pattern == nil {
			t.Errorf("rule %q has nil pattern", r.Name)
		}
	}
}
