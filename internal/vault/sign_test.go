package vault

import (
	"strings"
	"testing"
)

func TestSignResult_Fields(t *testing.T) {
	data := map[string]interface{}{"username": "admin", "password": "s3cr3t"}
	res := SignSecret("secret/app/creds", data, "my-signing-key")

	if !res.Success {
		t.Fatalf("expected Success=true, got false")
	}
	if res.Path != "secret/app/creds" {
		t.Errorf("expected path %q, got %q", "secret/app/creds", res.Path)
	}
	if res.Algorithm != "hmac-sha256" {
		t.Errorf("expected algorithm hmac-sha256, got %q", res.Algorithm)
	}
	if res.Signature == "" {
		t.Error("expected non-empty signature")
	}
	if res.KeyCount != 2 {
		t.Errorf("expected KeyCount=2, got %d", res.KeyCount)
	}
	if res.SignedAt.IsZero() {
		t.Error("expected non-zero SignedAt")
	}
	if res.Err != nil {
		t.Errorf("unexpected error: %v", res.Err)
	}
}

func TestSignResult_FailureState(t *testing.T) {
	res := SignSecret("secret/missing", nil, "key")

	if res.Success {
		t.Error("expected Success=false for nil data")
	}
	if res.Err == nil {
		t.Error("expected non-nil Err for nil data")
	}
	if !strings.Contains(res.Err.Error(), "nil data") {
		t.Errorf("error message should mention 'nil data', got: %v", res.Err)
	}
}

func TestSignSecret_Deterministic(t *testing.T) {
	data := map[string]interface{}{"b": "two", "a": "one", "c": "three"}
	res1 := SignSecret("p", data, "key")
	res2 := SignSecret("p", data, "key")

	if res1.Signature != res2.Signature {
		t.Errorf("signatures should be deterministic: %q vs %q", res1.Signature, res2.Signature)
	}
}

func TestSignSecret_DifferentKeysDifferentSig(t *testing.T) {
	data := map[string]interface{}{"k": "v"}
	res1 := SignSecret("p", data, "key-one")
	res2 := SignSecret("p", data, "key-two")

	if res1.Signature == res2.Signature {
		t.Error("different signing keys should produce different signatures")
	}
}

func TestVerifySecret_ValidSignature(t *testing.T) {
	data := map[string]interface{}{"env": "prod", "token": "abc123"}
	res := SignSecret("secret/svc", data, "verify-key")

	if !VerifySecret(data, "verify-key", res.Signature) {
		t.Error("expected VerifySecret to return true for valid signature")
	}
}

func TestVerifySecret_InvalidSignature(t *testing.T) {
	data := map[string]interface{}{"env": "prod"}
	if VerifySecret(data, "key", "deadbeef") {
		t.Error("expected VerifySecret to return false for wrong signature")
	}
}

func TestVerifySecret_NilData(t *testing.T) {
	if VerifySecret(nil, "key", "sig") {
		t.Error("expected VerifySecret to return false for nil data")
	}
}
