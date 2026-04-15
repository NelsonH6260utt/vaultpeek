package vault_test

import (
	"testing"

	"github.com/user/vaultpeek/internal/vault"
)

func TestKVv2DataPath(t *testing.T) {
	tests := []struct {
		mount    string
		path     string
		want     string
	}{
		{"secret", "myapp/config", "secret/data/myapp/config"},
		{"secret/", "/myapp/config", "secret/data/myapp/config"},
		{"kv", "prod/db", "kv/data/prod/db"},
	}

	for _, tt := range tests {
		t.Run(tt.mount+":"+tt.path, func(t *testing.T) {
			got := vault.KVv2DataPath(tt.mount, tt.path)
			if got != tt.want {
				t.Errorf("KVv2DataPath(%q, %q) = %q; want %q", tt.mount, tt.path, got, tt.want)
			}
		})
	}
}

func TestKVv2MetaPath(t *testing.T) {
	tests := []struct {
		mount string
		path  string
		want  string
	}{
		{"secret", "myapp", "secret/metadata/myapp"},
		{"secret/", "/myapp", "secret/metadata/myapp"},
	}

	for _, tt := range tests {
		t.Run(tt.mount+":"+tt.path, func(t *testing.T) {
			got := vault.KVv2MetaPath(tt.mount, tt.path)
			if got != tt.want {
				t.Errorf("KVv2MetaPath(%q, %q) = %q; want %q", tt.mount, tt.path, got, tt.want)
			}
		})
	}
}

func TestSecretData_TypeAssertion(t *testing.T) {
	data := vault.SecretData{
		"key1": "value1",
		"key2": 42,
	}

	if v, ok := data["key1"].(string); !ok || v != "value1" {
		t.Errorf("expected key1=%q, got %v", "value1", data["key1"])
	}
	if len(data) != 2 {
		t.Errorf("expected 2 keys, got %d", len(data))
	}
}
