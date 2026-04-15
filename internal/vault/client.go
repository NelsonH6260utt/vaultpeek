package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with vaultpeek-specific helpers.
type Client struct {
	api     *vaultapi.Client
	Address string
	Token   string
}

// Config holds the configuration needed to create a Vault client.
type Config struct {
	Address string
	Token   string
	CAPath  string
}

// NewClient creates and configures a new Vault client from the provided config.
// If Address or Token are empty, it falls back to environment variables.
func NewClient(cfg Config) (*Client, error) {
	vaultCfg := vaultapi.DefaultConfig()

	address := cfg.Address
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if address == "" {
		return nil, fmt.Errorf("vault address not set: provide --addr flag or VAULT_ADDR env var")
	}
	vaultCfg.Address = address

	if cfg.CAPath != "" {
		if err := vaultCfg.ConfigureTLS(&vaultapi.TLSConfig{CAPath: cfg.CAPath}); err != nil {
			return nil, fmt.Errorf("configuring TLS: %w", err)
		}
	}

	apiClient, err := vaultapi.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault api client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("vault token not set: provide --token flag or VAULT_TOKEN env var")
	}
	apiClient.SetToken(token)

	return &Client{
		api:     apiClient,
		Address: address,
		Token:   token,
	}, nil
}

// Logical returns the underlying logical client for raw Vault operations.
func (c *Client) Logical() *vaultapi.Logical {
	return c.api.Logical()
}
