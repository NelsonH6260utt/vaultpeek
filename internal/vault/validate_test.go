package vault

import (
	"testing"
)

func TestValidateSecret_Valid(t *testing.T) {
	data := map[string]interface{}{"username": "admin", "password": "s3cr3t"}
	result := ValidateSecret("secret/app", data, DefaultValidationRules())
	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
	if result.Checked != 2 {
		t.Errorf("expected 2 checked, got %d", result.Checked)
	}
}

func TestValidateSecret_EmptyValue(t *testing.T) {
	data := map[string]interface{}{"username": ""}
	result := ValidateSecret("secret/app", data, DefaultValidationRules())
	if result.Valid {
		t.Error("expected invalid due to empty value")
	}
	if len(result.Errors) == 0 {
		t.Error("expected at least one error")
	}
}

func TestValidateSecret_KeyWithSpace(t *testing.T) {
	data := map[string]interface{}{"bad key": "value"}
	result := ValidateSecret("secret/app", data, DefaultValidationRules())
	if result.Valid {
		t.Error("expected invalid due to key with space")
	}
}

func TestValidateSecret_NilData(t *testing.T) {
	result := ValidateSecret("secret/app", nil, DefaultValidationRules())
	if result.Valid {
		t.Error("expected invalid for nil data")
	}
	if result.Checked != 0 {
		t.Errorf("expected 0 checked, got %d", result.Checked)
	}
}

func TestValidateSecret_CustomRule(t *testing.T) {
	rule := func(key, value string) error {
		if len(value) < 8 {
			return nil
		}
		return nil
	}
	data := map[string]interface{}{"token": "abc"}
	result := ValidateSecret("secret/app", data, []ValidationRule{rule})
	if !result.Valid {
		t.Error("expected valid with permissive custom rule")
	}
}

func TestValidateSecret_PathPreserved(t *testing.T) {
	result := ValidateSecret("secret/myapp/prod", map[string]interface{}{"k": "v"}, DefaultValidationRules())
	if result.Path != "secret/myapp/prod" {
		t.Errorf("expected path preserved, got %q", result.Path)
	}
}
