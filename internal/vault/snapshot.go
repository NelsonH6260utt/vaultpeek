package vault

import (
	"fmt"
	"sort"
	"time"
)

// SnapshotEntry holds a single secret path and its data at snapshot time.
type SnapshotEntry struct {
	Path string
	Data map[string]interface{}
}

// SnapshotResult holds the result of a snapshot operation.
type SnapshotResult struct {
	TakenAt time.Time
	Entries []SnapshotEntry
	Success bool
	Err     error
}

// TakeSnapshot walks all paths under root and captures their secret data.
func TakeSnapshot(client *Client, mount, root string) SnapshotResult {
	result := SnapshotResult{
		TakenAt: time.Now().UTC(),
	}

	paths, err := WalkPaths(client, mount, root)
	if err != nil {
		result.Err = fmt.Errorf("walk failed: %w", err)
		return result
	}

	for _, p := range paths {
		data, err := FetchSecretData(client, mount, p)
		if err != nil {
			result.Err = fmt.Errorf("fetch failed for %s: %w", p, err)
			return result
		}
		result.Entries = append(result.Entries, SnapshotEntry{Path: p, Data: data})
	}

	sort.Slice(result.Entries, func(i, j int) bool {
		return result.Entries[i].Path < result.Entries[j].Path
	})

	result.Success = true
	return result
}
