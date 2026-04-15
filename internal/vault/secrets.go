package vault

import (
	"fmt"
	"strings"
)

// SecretData represents the key-value pairs stored at a Vault path.
type SecretData map[string]interface{}

// ListPaths returns the list of sub-paths (keys/directories) under the given path.
// It uses the KV v2 metadata list endpoint when mountPath is provided.
func (c *Client) ListPaths(path string) ([]string, error) {
	secret, err := c.api.Logical().List(path)
	if err != nil {
		return nil, fmt.Errorf("listing path %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	raw, ok := secret.Data["keys"]
	if !ok {
		return []string{}, nil
	}

	ifaces, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected keys type at path %q", path)
	}

	keys := make([]string, 0, len(ifaces))
	for _, v := range ifaces {
		if s, ok := v.(string); ok {
			keys = append(keys, s)
		}
	}
	return keys, nil
}

// ReadSecret reads the secret data at the given KV v2 path.
// path should be the full API path, e.g. "secret/data/myapp".
func (c *Client) ReadSecret(path string) (SecretData, error) {
	secret, err := c.api.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return SecretData{}, nil
	}

	// KV v2 wraps data under secret.Data["data"]
	if nested, ok := secret.Data["data"]; ok {
		if m, ok := nested.(map[string]interface{}); ok {
			return SecretData(m), nil
		}
	}

	return SecretData(secret.Data), nil
}

// KVv2DataPath converts a mount and secret path into the KV v2 data API path.
// e.g. mount="secret", path="myapp/config" -> "secret/data/myapp/config"
func KVv2DataPath(mount, path string) string {
	mount = strings.TrimSuffix(mount, "/")
	path = strings.TrimPrefix(path, "/")
	return fmt.Sprintf("%s/data/%s", mount, path)
}

// KVv2MetaPath converts a mount and secret path into the KV v2 metadata API path.
func KVv2MetaPath(mount, path string) string {
	mount = strings.TrimSuffix(mount, "/")
	path = strings.TrimPrefix(path, "/")
	return fmt.Sprintf("%s/metadata/%s", mount, path)
}
