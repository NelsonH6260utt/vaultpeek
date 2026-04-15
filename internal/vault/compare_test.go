package vault_test

import (
	"testing"

	"github.com/yourusername/vaultpeek/internal/vault"
)

func TestSecretMap_IsMapType(t *testing.T) {
	sm := vault.SecretMap{
		"username": "admin",
		"password": "s3cr3t",
	}

	if len(sm) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(sm))
	}

	if sm["username"] != "admin" {
		t.Errorf("expected username=admin, got %v", sm["username"])
	}
}

func TestSecretMap_EmptyMap(t *testing.T) {
	sm := vault.SecretMap{}
	if len(sm) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(sm))
	}
}

func TestSecretMap_NilSafe(t *testing.T) {
	var sm vault.SecretMap
	if sm != nil {
		t.Fatal("expected nil SecretMap to be nil")
	}
	// Accessing a nil map should not panic for reads.
	val := sm["key"]
	if val != nil {
		t.Errorf("expected nil for missing key, got %v", val)
	}
}

func TestSecretMap_Overwrite(t *testing.T) {
	sm := vault.SecretMap{
		"key": "original",
	}
	sm["key"] = "updated"

	if sm["key"] != "updated" {
		t.Errorf("expected updated value, got %v", sm["key"])
	}
}
