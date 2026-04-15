package vault

import (
	"context"
	"fmt"
	"sort"
	"strconv"
)

// VersionMeta holds metadata for a single version of a KV v2 secret.
type VersionMeta struct {
	Version      int
	CreatedTime  string
	DeletionTime string
	Destroyed    bool
}

// ListVersions returns version metadata for all versions of a KV v2 secret.
func ListVersions(ctx context.Context, c *Client, mount, secretPath string) ([]VersionMeta, error) {
	metaPath := KVv2MetaPath(mount, secretPath)

	secret, err := c.Logical().ReadWithContext(ctx, metaPath)
	if err != nil {
		return nil, fmt.Errorf("reading metadata at %s: %w", metaPath, err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no metadata found at %s", metaPath)
	}

	versionsRaw, ok := secret.Data["versions"]
	if !ok {
		return nil, fmt.Errorf("no versions key in metadata at %s", metaPath)
	}

	versionsMap, ok := versionsRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for versions field")
	}

	var metas []VersionMeta
	for vStr, vRaw := range versionsMap {
		vNum, err := strconv.Atoi(vStr)
		if err != nil {
			continue
		}
		vMap, ok := vRaw.(map[string]interface{});
		if !ok {
			continue
		}
		meta := VersionMeta{Version: vNum}
		if ct, ok := vMap["created_time"].(string); ok {
			meta.CreatedTime = ct
		}
		if dt, ok := vMap["deletion_time"].(string); ok {
			meta.DeletionTime = dt
		}
		if d, ok := vMap["destroyed"].(bool); ok {
			meta.Destroyed = d
		}
		metas = append(metas, meta)
	}

	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Version < metas[j].Version
	})
	return metas, nil
}
