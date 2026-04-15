package vault

import (
	"context"
	"fmt"
	"strings"
)

// ListSecrets lists all secret keys under the given path in a KV v2 mount.
// It uses the metadata endpoint to enumerate keys.
func (c *Client) ListSecrets(ctx context.Context, mount, path string) ([]string, error) {
	metaPath := KVv2MetaPath(mount, path)

	secret, err := c.vault.Logical().ListWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("listing secrets at %q: %w", metaPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no secrets found at path %q", metaPath)
	}

	raw, ok := secret.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("unexpected response: missing 'keys' field at %q", metaPath)
	}

	ifaces, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for 'keys' field at %q", metaPath)
	}

	keys := make([]string, 0, len(ifaces))
	for _, v := range ifaces {
		s, ok := v.(string)
		if !ok {
			continue
		}
		keys = append(keys, s)
	}
	return keys, nil
}

// IsDir reports whether a key returned from ListSecrets represents a
// sub-directory (Vault appends a trailing slash to folder entries).
func IsDir(key string) bool {
	return strings.HasSuffix(key, "/")
}

// TrimDir removes the trailing slash from a directory key.
func TrimDir(key string) string {
	return strings.TrimSuffix(key, "/")
}
