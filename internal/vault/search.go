package vault

import (
	"context"
	"strings"
)

// SearchResult holds a matching secret path and the key that matched.
type SearchResult struct {
	Path string
	MatchedKey string
}

// SearchSecrets recursively walks all paths under root and returns paths
// whose secret keys or values contain the given query string (case-insensitive).
func SearchSecrets(ctx context.Context, c *Client, mount, root, query string) ([]SearchResult, error) {
	paths, err := WalkPaths(ctx, c, mount, root)
	if err != nil {
		return nil, err
	}

	q := strings.ToLower(query)
	var results []SearchResult

	for _, p := range paths {
		dataPath := KVv2DataPath(mount, p)
		secret, err := c.Logical().ReadWithContext(ctx, dataPath)
		if err != nil || secret == nil || secret.Data == nil {
			continue
		}

		data, ok := secret.Data["data"].(map[string]interface{})
		if !ok {
			continue
		}

		for k, v := range data {
			if strings.Contains(strings.ToLower(k), q) {
				results = append(results, SearchResult{Path: p, MatchedKey: k})
				break
			}
			if vs, ok := v.(string); ok && strings.Contains(strings.ToLower(vs), q) {
				results = append(results, SearchResult{Path: p, MatchedKey: k})
				break
			}
		}
	}

	return results, nil
}

// WalkPaths recursively lists all leaf secret paths under root.
func WalkPaths(ctx context.Context, c *Client, mount, root string) ([]string, error) {
	entries, err := ListSecrets(ctx, c, mount, root)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, e := range entries {
		full := root + TrimDir(e)
		if IsDir(e) {
			sub, err := WalkPaths(ctx, c, mount, full+"/")
			if err != nil {
				return nil, err
			}
			paths = append(paths, sub...)
		} else {
			paths = append(paths, full)
		}
	}
	return paths, nil
}
