package vault

import (
	"fmt"
	"time"
)

// RotateResult holds the outcome of a secret rotation operation.
type RotateResult struct {
	Path       string
	OldVersion int
	NewVersion int
	RotatedAt  time.Time
	Keys       []string
	Success    bool
	Err        error
}

// RotateOptions configures the rotation behaviour.
type RotateOptions struct {
	// Transform is an optional function applied to each value during rotation.
	// If nil, values are carried forward unchanged.
	Transform func(key, value string) string
}

// RotateSecret reads the current data at path, optionally transforms each
// value, and writes a new version back to Vault. This is useful for
// re-encrypting or refreshing secrets without changing their keys.
func RotateSecret(c *Client, path string, opts RotateOptions) RotateResult {
	result := RotateResult{Path: path, RotatedAt: time.Now().UTC()}

	secret, err := c.Logical().Read(KVv2DataPath(c.Mount, path))
	if err != nil {
		result.Err = fmt.Errorf("rotate: read %s: %w", path, err)
		return result
	}
	if secret == nil || secret.Data == nil {
		result.Err = fmt.Errorf("rotate: no data at %s", path)
		return result
	}

	raw, err := SecretData(secret)
	if err != nil {
		result.Err = fmt.Errorf("rotate: parse data %s: %w", path, err)
		return result
	}

	if meta, ok := secret.Data["metadata"].(map[string]interface{}); ok {
		if v, ok := meta["version"].(float64); ok {
			result.OldVersion = int(v)
		}
	}

	newData := make(map[string]interface{}, len(raw))
	for k, v := range raw {
		strVal := fmt.Sprintf("%v", v)
		if opts.Transform != nil {
			strVal = opts.Transform(k, strVal)
		}
		newData[k] = strVal
		result.Keys = append(result.Keys, k)
	}

	written, err := c.Logical().Write(KVv2DataPath(c.Mount, path), map[string]interface{}{"data": newData})
	if err != nil {
		result.Err = fmt.Errorf("rotate: write %s: %w", path, err)
		return result
	}

	if written != nil {
		if v, ok := written.Data["version"].(float64); ok {
			result.NewVersion = int(v)
		}
	}

	result.Success = true
	return result
}
