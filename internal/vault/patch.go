package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// PatchResult holds the outcome of a partial secret update operation.
type PatchResult struct {
	Path        string
	UpdatedKeys []string
	Success     bool
	Error       error
}

// PatchSecret performs a partial update on a KV v2 secret at the given path,
// merging the provided fields into the existing secret data without overwriting
// keys that are not present in the patch map.
func PatchSecret(ctx context.Context, client *api.Client, mount, path string, patch map[string]interface{}) PatchResult {
	fullPath := KVv2DataPath(mount, path)

	// Read the existing secret first.
	existing, err := client.Logical().ReadWithContext(ctx, fullPath)
	if err != nil {
		return PatchResult{Path: path, Success: false, Error: fmt.Errorf("read existing secret: %w", err)}
	}

	merged := make(map[string]interface{})
	if existing != nil && existing.Data != nil {
		if data, ok := existing.Data["data"].(map[string]interface{}); ok {
			for k, v := range data {
				merged[k] = v
			}
		}
	}

	updatedKeys := make([]string, 0, len(patch))
	for k, v := range patch {
		merged[k] = v
		updatedKeys = append(updatedKeys, k)
	}

	_, err = client.Logical().WriteWithContext(ctx, fullPath, map[string]interface{}{
		"data": merged,
	})
	if err != nil {
		return PatchResult{Path: path, Success: false, Error: fmt.Errorf("write patched secret: %w", err)}
	}

	return PatchResult{
		Path:        path,
		UpdatedKeys: updatedKeys,
		Success:     true,
	}
}
