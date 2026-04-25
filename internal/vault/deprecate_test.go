package vault

import (
	"testing"
	"time"
)

func TestDeprecateResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	r := DeprecateResult{
		Path:       "secret/myapp/db",
		Succeeded:  true,
		Reason:     "use secret/myapp/database instead",
		Deprecated: now,
	}
	if r.Path != "secret/myapp/db" {
		t.Errorf("expected path to be set, got %q", r.Path)
	}
	if !r.Succeeded {
		t.Error("expected Succeeded to be true")
	}
	if r.Reason == "" {
		t.Error("expected Reason to be set")
	}
	if r.Deprecated.IsZero() {
		t.Error("expected Deprecated timestamp to be set")
	}
}

func TestDeprecateResult_FailureState(t *testing.T) {
	r := DeprecateResult{
		Path:      "secret/old",
		Succeeded: false,
		Error:     fmt.Errorf("vault unavailable"),
	}
	if r.Succeeded {
		t.Error("expected Succeeded to be false")
	}
	if r.Error == nil {
		t.Error("expected Error to be set")
	}
}

func TestIsDeprecated_True(t *testing.T) {
	meta := map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"deprecated": "true",
		},
	}
	if !IsDeprecated(meta) {
		t.Error("expected IsDeprecated to return true")
	}
}

func TestIsDeprecated_False(t *testing.T) {
	meta := map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"deprecated": "false",
		},
	}
	if IsDeprecated(meta) {
		t.Error("expected IsDeprecated to return false")
	}
}

func TestIsDeprecated_Missing(t *testing.T) {
	if IsDeprecated(map[string]interface{}{}) {
		t.Error("expected IsDeprecated to return false for empty metadata")
	}
}

func TestIsDeprecated_NilCustomMetadata(t *testing.T) {
	meta := map[string]interface{}{
		"custom_metadata": nil,
	}
	if IsDeprecated(meta) {
		t.Error("expected IsDeprecated to return false when custom_metadata is nil")
	}
}
