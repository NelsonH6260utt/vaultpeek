package vault

import (
	"fmt"
	"sort"
)

// TagResult holds the outcome of a tag operation on a secret path.
type TagResult struct {
	Path    string
	Tags    map[string]string
	Success bool
	Err     error
}

// TagSecret writes custom_metadata tags to a KVv2 secret path.
func TagSecret(client *Client, mount, path string, tags map[string]string) TagResult {
	metaPath := KVv2MetaPath(mount, path)

	existing, err := client.Logical().Read(metaPath)
	if err != nil {
		return TagResult{Path: path, Err: fmt.Errorf("read metadata: %w", err)}
	}

	merged := map[string]interface{}{}
	if existing != nil && existing.Data != nil {
		if cm, ok := existing.Data["custom_metadata"].(map[string]interface{}); ok {
			for k, v := range cm {
				merged[k] = v
			}
		}
	}
	for k, v := range tags {
		merged[k] = v
	}

	_, err = client.Logical().Write(metaPath, map[string]interface{}{
		"custom_metadata": merged,
	})
	if err != nil {
		return TagResult{Path: path, Err: fmt.Errorf("write metadata: %w", err)}
	}

	result := make(map[string]string, len(merged))
	for k, v := range merged {
		result[k] = fmt.Sprintf("%v", v)
	}
	return TagResult{Path: path, Tags: result, Success: true}
}

// SortedTagKeys returns tag keys in sorted order.
func SortedTagKeys(tags map[string]string) []string {
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
