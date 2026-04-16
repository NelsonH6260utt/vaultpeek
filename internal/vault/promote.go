package vault

import (
	"fmt"
	"strings"
)

// PromoteResult holds the outcome of promoting a secret from one environment path to another.
type PromoteResult struct {
	SourcePath string
	DestPath   string
	Keys       []string
	Success    bool
	Err        error
}

// PromoteSecret copies a secret from a source path to a destination path,
// typically used to promote secrets across environments (e.g., staging -> production).
func PromoteSecret(client *Client, mount, srcPath, destPath string) PromoteResult {
	result := PromoteResult{
		SourcePath: srcPath,
		DestPath:   destPath,
	}

	if strings.TrimSpace(srcPath) == "" || strings.TrimSpace(destPath) == "" {
		result.Err = fmt.Errorf("source and destination paths must not be empty")
		return result
	}

	data, err := FetchSecretData(client, mount, srcPath)
	if err != nil {
		result.Err = fmt.Errorf("failed to read source secret: %w", err)
		return result
	}

	dataPath := KVv2DataPath(mount, destPath)
	_, err = client.Logical().Write(dataPath, map[string]interface{}{"data": data})
	if err != nil {
		result.Err = fmt.Errorf("failed to write destination secret: %w", err)
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
