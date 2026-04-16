package vault

import (
	"fmt"
	"time"
)

// RestoreResult holds the outcome of a secret restore operation.
type RestoreResult struct {
	Path       string
	FromVersion int
	NewVersion  int
	RestoredAt  time.Time
	Keys        []string
	Success     bool
	Err         error
}

// RestoreSecret restores a KV v2 secret to a specific historical version
// by reading that version and writing it as a new version.
func RestoreSecret(client *Client, mount, path string, version int) RestoreResult {
	dataPath := KVv2DataPath(mount, path)

	// Read the target version
	params := map[string][]string{"version": {fmt.Sprintf("%d", version)}}
	secret, err := client.Logical().ReadWithData(dataPath, params)
	if err != nil {
		return RestoreResult{Path: path, Success: false, Err: err}
	}
	if secret == nil || secret.Data == nil {
		return RestoreResult{Path: path, Success: false, Err: fmt.Errorf("version %d not found at %s", version, path)}
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return RestoreResult{Path: path, Success: false, Err: fmt.Errorf("unexpected data format at %s", path)}
	}

	// Write restored data as a new version
	writeSecret, err := client.Logical().Write(dataPath, map[string]interface{}{"data": data})
	if err != nil {
		return RestoreResult{Path: path, Success: false, Err: err}
	}

	newVersion := 0
	if writeSecret != nil {
		if meta, ok := writeSecret.Data["version"].(json.Number); ok {
			newVersion, _ = strconv.Atoi(meta.String())
		}
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	return RestoreResult{
		Path:        path,
		FromVersion: version,
		NewVersion:  newVersion,
		RestoredAt:  time.Now().UTC(),
		Keys:        keys,
		Success:     true,
	}
}
