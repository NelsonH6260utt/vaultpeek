package vault

import (
	"fmt"
	"strings"
)

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	SourcePath string
	DestPath   string
	Keys       []string
	Success    bool
	Error      error
}

// CloneSecret reads all key-value pairs from sourcePath and writes them to
// destPath as a new KV v2 secret. It does not overwrite existing keys at
// destPath — callers should check before invoking if that matters.
func CloneSecret(c *Client, mount, sourcePath, destPath string) CloneResult {
	result := CloneResult{
		SourcePath: sourcePath,
		DestPath:   destPath,
	}

	if strings.TrimSpace(sourcePath) == "" || strings.TrimSpace(destPath) == "" {
		result.Error = fmt.Errorf("source and destination paths must not be empty")
		return result
	}

	dataPath := KVv2DataPath(mount, sourcePath)
	secret, err := c.Logical().Read(dataPath)
	if err != nil {
		result.Error = fmt.Errorf("read source: %w", err)
		return result
	}
	if secret == nil || secret.Data == nil {
		result.Error = fmt.Errorf("source path %q not found", sourcePath)
		return result
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		result.Error = fmt.Errorf("unexpected data format at %q", sourcePath)
		return result
	}

	destDataPath := KVv2DataPath(mount, destPath)
	_, err = c.Logical().Write(destDataPath, map[string]interface{}{"data": data})
	if err != nil {
		result.Error = fmt.Errorf("write dest: %w", err)
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
