package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// RollbackResult holds the outcome of a secret version rollback.
type RollbackResult struct {
	Path        string
	FromVersion int
	ToVersion   int
	Success     bool
	Error       error
}

// RollbackSecret rolls back a KV v2 secret at the given path to the specified version.
// It reads the target version's data and writes it as a new version.
func RollbackSecret(ctx context.Context, client *api.Client, mount, path string, toVersion int) RollbackResult {
	fullDataPath := KVv2DataPath(mount, path)
	fullMetaPath := KVv2MetaPath(mount, path)

	// Read current metadata to determine the current version number
	meta, err := client.Logical().ReadWithContext(ctx, fullMetaPath)
	if err != nil {
		return RollbackResult{Path: path, Error: fmt.Errorf("reading metadata: %w", err)}
	}
	if meta == nil {
		return RollbackResult{Path: path, Error: fmt.Errorf("secret not found: %s", path)}
	}

	currentVersion := 0
	if cv, ok := meta.Data["current_version"]; ok {
		if n, ok := cv.(float64); ok {
			currentVersion = int(n)
		}
	}

	// Read the target version's data
	targetData, err := client.Logical().ReadWithDataWithContext(ctx, fullDataPath, map[string][]string{
		"version": {fmt.Sprintf("%d", toVersion)},
	})
	if err != nil {
		return RollbackResult{Path: path, Error: fmt.Errorf("reading version %d: %w", toVersion, err)}
	}
	if targetData == nil {
		return RollbackResult{Path: path, Error: fmt.Errorf("version %d not found at %s", toVersion, path)}
	}

	secretData, ok := targetData.Data["data"].(map[string]interface{})
	if !ok {
		return RollbackResult{Path: path, Error: fmt.Errorf("unexpected data format at version %d", toVersion)}
	}

	// Write the old version's data as a new version
	_, err = client.Logical().WriteWithContext(ctx, fullDataPath, map[string]interface{}{
		"data": secretData,
	})
	if err != nil {
		return RollbackResult{Path: path, Error: fmt.Errorf("writing rollback data: %w", err)}
	}

	return RollbackResult{
		Path:        path,
		FromVersion: currentVersion,
		ToVersion:   toVersion,
		Success:     true,
	}
}
