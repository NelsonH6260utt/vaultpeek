package vault

import (
	"fmt"
	"time"
)

// ArchiveEntry holds a single archived secret path and its data.
type ArchiveEntry struct {
	Path      string
	Data      map[string]interface{}
	ArchivedAt time.Time
}

// ArchiveResult is returned by ArchiveSecret.
type ArchiveResult struct {
	Path    string
	Success bool
	Error   error
	Keys    []string
}

// ArchiveSecret reads a secret at path and returns an ArchiveEntry along with
// a result summary. It does not delete the original secret.
func ArchiveSecret(client *Client, mount, path string) (ArchiveEntry, ArchiveResult) {
	result := ArchiveResult{Path: path}

	data, err := FetchSecretData(client, mount, path)
	if err != nil {
		result.Error = fmt.Errorf("archive: read %q: %w", path, err)
		return ArchiveEntry{}, result
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}

	result.Success = true
	result.Keys = keys

	return ArchiveEntry{
		Path:       path,
		Data:       data,
		ArchivedAt: time.Now().UTC(),
	}, result
}
