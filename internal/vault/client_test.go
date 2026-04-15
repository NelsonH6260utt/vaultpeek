package vault_test

import (
	"os"
	"testing"

	"github.com/user/vaultpeek/internal/vault"
)

func TestNewClient_MissingAddress(t *testing.T) {
	os.Unsetenv("VAULT_ADDR")
	os.Unsetenv("VAULT_TOKEN")

	_, err := vault.NewClient(vault.Config{})
	if err == nil {
		t.Fatal("expected error when address is missing, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	os.Unsetenv("VAULT_TOKEN")

	_, err := vault.NewClient(vault.Config{
		Address: "http://127.0.0.1:8200",
	})
	if err == nil {
		t.Fatal("expected error when token is missing, got nil")
	}
}

func TestNewClient_FromEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "test-token")

	c, err := vault.NewClient(vault.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Address != "http://127.0.0.1:8200" {
		t.Errorf("expected address %q, got %q", "http://127.0.0.1:8200", c.Address)
	}
}

func TestNewClient_ExplicitConfig(t *testing.T) {
	c, err := vault.NewClient(vault.Config{
		Address: "http://vault.example.com:8200",
		Token:   "s.mytoken",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Token != "s.mytoken" {
		t.Errorf("expected token %q, got %q", "s.mytoken", c.Token)
	}
}
