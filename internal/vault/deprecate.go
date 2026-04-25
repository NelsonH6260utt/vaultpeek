package vault

import (
	"fmt"
	"time"
)

// DeprecateResult holds the outcome of marking a secret path as deprecated.
type DeprecateResult struct {
	Path       string
	Succeeded  bool
	Reason     string
	Deprecated time.Time
	Error      error
}

// DeprecateOptions controls how a secret is marked deprecated.
type DeprecateOptions struct {
	Reason      string
	Replacement string
}

// DeprecateSecret marks a secret path as deprecated by writing metadata
// annotations indicating the deprecation reason and timestamp.
func DeprecateSecret(c *Client, path string, opts DeprecateOptions) DeprecateResult {
	if path == "" {
		return DeprecateResult{
			Path:      path,
			Succeeded: false,
			Error:     fmt.Errorf("path must not be empty"),
		}
	}

	now := time.Now().UTC()

	annotations := map[string]string{
		"deprecated":        "true",
		"deprecated-at":     now.Format(time.RFC3339),
		"deprecated-reason": opts.Reason,
	}
	if opts.Replacement != "" {
		annotations["deprecated-replacement"] = opts.Replacement
	}

	metaPath := KVv2MetaPath(path)
	_, err := c.Logical().Write(metaPath, map[string]interface{}{
		"custom_metadata": annotations,
	})
	if err != nil {
		return DeprecateResult{
			Path:      path,
			Succeeded: false,
			Reason:    opts.Reason,
			Error:     fmt.Errorf("writing deprecation metadata: %w", err),
		}
	}

	return DeprecateResult{
		Path:       path,
		Succeeded:  true,
		Reason:     opts.Reason,
		Deprecated: now,
	}
}

// IsDeprecated reports whether a secret's custom metadata contains a
// deprecation marker.
func IsDeprecated(meta map[string]interface{}) bool {
	cm, ok := meta["custom_metadata"]
	if !ok || cm == nil {
		return false
	}
	fields, ok := cm.(map[string]interface{})
	if !ok {
		return false
	}
	v, ok := fields["deprecated"]
	if !ok {
		return false
	}
	s, ok := v.(string)
	return ok && s == "true"
}
