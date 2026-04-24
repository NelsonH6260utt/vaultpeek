package vault

import (
	"fmt"
	"sort"
)

// LabelResult holds the outcome of a label operation on a Vault secret.
type LabelResult struct {
	Path    string
	Labels  map[string]string
	Success bool
	Err     error
}

// SortedLabelKeys returns the keys of a label map in sorted order.
func SortedLabelKeys(labels map[string]string) []string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// LabelSecret writes label metadata to a KVv2 secret's custom_metadata.
// Labels are stored under a "labels/" prefix to avoid collisions.
func LabelSecret(client *Client, path string, labels map[string]string) LabelResult {
	if client == nil {
		return LabelResult{
			Path:    path,
			Success: false,
			Err:     fmt.Errorf("nil client"),
		}
	}

	if len(labels) == 0 {
		return LabelResult{
			Path:    path,
			Labels:  map[string]string{},
			Success: true,
		}
	}

	customMeta := make(map[string]string, len(labels))
	for k, v := range labels {
		customMeta["label/"+k] = v
	}

	metaPath := KVv2MetaPath(path)
	_, err := client.Logical().Write(metaPath, map[string]interface{}{
		"custom_metadata": customMeta,
	})
	if err != nil {
		return LabelResult{
			Path:    path,
			Success: false,
			Err:     fmt.Errorf("writing labels: %w", err),
		}
	}

	return LabelResult{
		Path:    path,
		Labels:  labels,
		Success: true,
	}
}
