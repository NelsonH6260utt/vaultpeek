package vault

import (
	"context"
	"time"
)

// WatchResult holds the outcome of a single watch poll cycle.
type WatchResult struct {
	Path    string
	Changed bool
	Version int
	Data    map[string]interface{}
	Error   error
}

// WatchOptions configures the watch behaviour.
type WatchOptions struct {
	Interval time.Duration
	MaxPolls int // 0 means unlimited
}

// WatchSecret polls a KV v2 secret at the given path and emits a WatchResult
// on the returned channel each time the version changes (or on error).
// The channel is closed when ctx is cancelled or MaxPolls is reached.
func WatchSecret(ctx context.Context, client *Client, mount, path string, opts WatchOptions) <-chan WatchResult {
	ch := make(chan WatchResult, 1)
	if opts.Interval <= 0 {
		opts.Interval = 5 * time.Second
	}

	go func() {
		defer close(ch)
		lastVersion := -1
		polls := 0

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(opts.Interval):
			}

			data, version, err := fetchVersioned(client, mount, path)
			result := WatchResult{Path: path, Version: version, Data: data, Error: err}
			if err == nil && version != lastVersion {
				result.Changed = true
				lastVersion = version
			}

			select {
			case ch <- result:
			case <-ctx.Done():
				return
			}

			polls++
			if opts.MaxPolls > 0 && polls >= opts.MaxPolls {
				return
			}
		}
	}()

	return ch
}

// fetchVersioned reads the secret and returns its data plus current version.
func fetchVersioned(client *Client, mount, path string) (map[string]interface{}, int, error) {
	dataPath := KVv2DataPath(mount, path)
	secret, err := client.Logical().Read(dataPath)
	if err != nil {
		return nil, 0, err
	}
	if secret == nil || secret.Data == nil {
		return map[string]interface{}{}, 0, nil
	}
	meta, _ := secret.Data["metadata"].(map[string]interface{})
	version := 0
	if meta != nil {
		if v, ok := meta["version"].(float64); ok {
			version = int(v)
		}
	}
	data, _ := secret.Data["data"].(map[string]interface{})
	if data == nil {
		data = map[string]interface{}{}
	}
	return data, version, nil
}
