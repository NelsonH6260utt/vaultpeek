package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// CopyResult holds the outcome of a secret copy operation.
type CopyResult struct {
	SourcePath string
	DestPath   string
	Keys       []string
	Success    bool
	Err        error
}

// CopySecret reads a KVv2 secret from src and writes it to dst using the
// provided mount. Both paths should be relative to the mount point.
func CopySecret(ctx context.Context, client *api.Client, mount, src, dst string) CopyResult {
	result := CopyResult{
		SourcePath: src,
		DestPath:   dst,
	}

	dataPath := KVv2DataPath(mount, src)
	secret, err := client.Logical().ReadWithContext(ctx, dataPath)
	if err != nil {
		result.Err = fmt.Errorf("reading source %q: %w", src, err)
		return result
	}
	if secret == nil || secret.Data == nil {
		result.Err = fmt.Errorf("source secret %q not found", src)
		return result
	}

	rawData, ok := secret.Data["data"]
	if !ok {
		result.Err = fmt.Errorf("source secret %q has no data field", src)
		return result
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		result.Err = fmt.Errorf("source secret %q data is not a map", src)
		return result
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	result.Keys = keys

	destPath := KVv2DataPath(mount, dst)
	_, err = client.Logical().WriteWithContext(ctx, destPath, map[string]interface{}{
		"data": data,
	})
	if err != nil {
		result.Err = fmt.Errorf("writing destination %q: %w", dst, err)
		return result
	}

	result.Success = true
	return result
}
