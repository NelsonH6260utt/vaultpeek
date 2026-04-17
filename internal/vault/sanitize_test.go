package vault

import (
	"testing"
)

func TestSanitizeResult_Fields(t *testing.T) {
	r := SanitizeResult{
		Path:    "secret/app",
		Success: true,
		Changed: map[string]string{"password": "[REDACTED]"},
		Keys:    []string{"password", "user"},
	}
	if r.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %s", r.Path)
	}
	if !r.Success {
		t.Error("expected success")
	}
	if r.Changed["password"] != "[REDACTED]" {
		t.Error("expected password to be redacted")
	}
}

func TestSanitizeResult_FailureState(t *testing.T) {
	r := SanitizeResult{Success: false, Error: errSentinel("sanitize failed")}
	if r.Success {
		t.Error("expected failure")
	}
	if r.Error == nil {
		t.Error("expected error to be set")
	}
}

func TestSanitizeSecret_NilData(t *testing.T) {
	r := SanitizeSecret("secret/app", nil, DefaultSanitizeRules())
	if !r.Success {
		t.Error("expected success on nil data")
	}
	if len(r.Changed) != 0 {
		t.Error("expected no changed keys")
	}
}

func TestSanitizeSecret_RedactsPassword(t *testing.T) {
	data := map[string]interface{}{
		"db_password": "supersecret",
		"username":    "admin",
	}
	r := SanitizeSecret("secret/db", data, DefaultSanitizeRules())
	if _, ok := r.Changed["db_password"]; !ok {
		t.Error("expected db_password to be flagged")
	}
	if _, ok := r.Changed["username"]; ok {
		t.Error("expected username not to be flagged")
	}
}

func TestSanitizeSecret_NoMatch(t *testing.T) {
	data := map[string]interface{}{"region": "us-east-1", "env": "prod"}
	r := SanitizeSecret("secret/cfg", data, DefaultSanitizeRules())
	if len(r.Changed) != 0 {
		t.Errorf("expected no redactions, got %d", len(r.Changed))
	}
}

func TestSanitizeSecret_CustomRule(t *testing.T) {
	rule := func(key, _ string) bool { return key == "pin" }
	data := map[string]interface{}{"pin": "1234", "name": "alice"}
	r := SanitizeSecret("secret/user", data, []SanitizeRule{rule})
	if r.Changed["pin"] != "[REDACTED]" {
		t.Error("expected pin to be redacted by custom rule")
	}
}

type errSentinel string

func (e errSentinel) Error() string { return string(e) }
