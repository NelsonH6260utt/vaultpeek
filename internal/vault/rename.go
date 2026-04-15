package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// RenameResult holds the outcome of a secret rename (copy + delete) operation.
type RenameResult struct {
	SourcePath string
	DestPath   string
	Keys       []string
	Success    bool
	Err        error
}

// RenameSecret copies a KVv2 secret from srcPath to dstPath, then deletes the
// source. Both paths are relative to the mount (e.g. "myapp/config").
func RenameSecret(ctx context.Context, client *api.Client, mount, srcPath, dstPath string) RenameResult {
	result := RenameResult{
		SourcePath: srcPath,
		DestPath:   dstPath,
	}

	// Read source secret.
	readPath := KVv2DataPath(mount, srcPath)
	secret, err := client.Logical().ReadWithContext(ctx, readPath)
	if err != nil {
		result.Err = fmt.Errorf("read source %q: %w", srcPath, err)
		return result
	}
	if secret == nil || secret.Data == nil {
		result.Err = fmt.Errorf("source secret %q not found", srcPath)
		return result
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		result.Err = fmt.Errorf("unexpected data format at %q", srcPath)
		return result
	}

	// Write to destination.
	writePath := KVv2DataPath(mount, dstPath)
	_, err = client.Logical().WriteWithContext(ctx, writePath, map[string]interface{}{"data": data})
	if err != nil {
		result.Err = fmt.Errorf("write dest %q: %w", dstPath, err)
		return result
	}

	// Delete the source (metadata delete removes all versions).
	metaPath := KVv2MetaPath(mount, srcPath)
	_, err = client.Logical().DeleteWithContext(ctx, metaPath)
	if err != nil {
		result.Err = fmt.Errorf("delete source %q: %w", srcPath, err)
		return result
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	result.Keys = keys
	result.Success = true
	return result
}
