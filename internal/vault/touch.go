package vault

import (
	"fmt"
	"time"
)

// TouchResult holds the outcome of a touch operation on a Vault secret.
type TouchResult struct {
	Path      string
	Success   bool
	Error     error
	TouchedAt time.Time
	Version   int
}

// TouchSecret re-writes a secret at the given path with its current data,
// effectively bumping its version and updating its last-modified timestamp.
// This is useful for triggering watchers or refreshing TTL-based metadata.
func TouchSecret(client *Client, mountPath, secretPath string) TouchResult {
	fullPath := KVv2DataPath(mountPath, secretPath)

	// Read current secret data
	secret, err := client.Logical().Read(fullPath)
	if err != nil {
		return TouchResult{
			Path:    secretPath,
			Success: false,
			Error:   fmt.Errorf("read failed: %w", err),
		}
	}
	if secret == nil || secret.Data == nil {
		return TouchResult{
			Path:    secretPath,
			Success: false,
			Error:   fmt.Errorf("secret not found at path: %s", secretPath),
		}
	}

	// Extract the inner data map from KV v2 response
	rawData, ok := secret.Data["data"]
	if !ok {
		return TouchResult{
			Path:    secretPath,
			Success: false,
			Error:   fmt.Errorf("unexpected secret format at path: %s", secretPath),
		}
	}
	dataMap, ok := rawData.(map[string]interface{})
	if !ok {
		return TouchResult{
			Path:    secretPath,
			Success: false,
			Error:   fmt.Errorf("data field is not a map at path: %s", secretPath),
		}
	}

	// Re-write the secret to bump its version
	payload := map[string]interface{}{
		"data": dataMap,
	}
	written, err := client.Logical().Write(fullPath, payload)
	if err != nil {
		return TouchResult{
			Path:    secretPath,
			Success: false,
			Error:   fmt.Errorf("write failed: %w", err),
		}
	}

	newVersion := 0
	if written != nil && written.Data != nil {
		if meta, ok := written.Data["version"]; ok {
			if v, ok := meta.(int); ok {
				newVersion = v
			}
		}
	}

	return TouchResult{
		Path:      secretPath,
		Success:   true,
		TouchedAt: time.Now().UTC(),
		Version:   newVersion,
	}
}
