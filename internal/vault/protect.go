package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// ProtectResult holds the outcome of a protect (write-protect) operation.
type ProtectResult struct {
	Path    string
	Mount   string
	Success bool
	Error   error
}

// ProtectSecret marks a KV v2 secret path as non-deletable by writing
// custom metadata with a "protected" flag.
func ProtectSecret(client *api.Client, mount, path string) ProtectResult {
	result := ProtectResult{
		Path:  path,
		Mount: mount,
	}

	metaPath := KVv2MetaPath(mount, path)

	_, err := client.Logical().WriteWithContext(context.Background(), metaPath, map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"vaultpeek_protected": "true",
		},
	})
	if err != nil {
		result.Error = fmt.Errorf("protect: write metadata for %q: %w", path, err)
		return result
	}

	result.Success = true
	return result
}

// IsProtected checks whether a KV v2 secret path has the vaultpeek protection
// flag set in its custom metadata.
func IsProtected(client *api.Client, mount, path string) (bool, error) {
	metaPath := KVv2MetaPath(mount, path)

	secret, err := client.Logical().ReadWithContext(context.Background(), metaPath)
	if err != nil {
		return false, fmt.Errorf("protect: read metadata for %q: %w", path, err)
	}
	if secret == nil || secret.Data == nil {
		return false, nil
	}

	cm, ok := secret.Data["custom_metadata"]
	if !ok {
		return false, nil
	}

	meta, ok := cm.(map[string]interface{})
	if !ok {
		return false, nil
	}

	val, _ := meta["vaultpeek_protected"].(string)
	return val == "true", nil
}
