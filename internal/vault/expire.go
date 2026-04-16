package vault

import (
	"fmt"
	"time"
)

// ExpireResult holds the outcome of setting a TTL/expiry on a secret.
type ExpireResult struct {
	Path      string
	TTL       time.Duration
	ExpiresAt time.Time
	Success   bool
	Err       error
}

// ExpireOptions configures expiry behaviour.
type ExpireOptions struct {
	TTL time.Duration
}

// SetExpiry writes an expiry timestamp into the secret's metadata custom_metadata.
// Vault OSS does not natively support TTL on KV secrets, so we store it as a tag.
func SetExpiry(c *Client, mount, path string, opts ExpireOptions) ExpireResult {
	if opts.TTL <= 0 {
		return ExpireResult{
			Path:    path,
			Success: false,
			Err:     fmt.Errorf("TTL must be greater than zero"),
		}
	}

	expiresAt := time.Now().UTC().Add(opts.TTL).Truncate(time.Second)

	metaPath := KVv2MetaPath(mount, path)
	body := map[string]interface{}{
		"custom_metadata": map[string]interface{}{
			"expires_at": expiresAt.Format(time.RFC3339),
			"ttl_seconds": fmt.Sprintf("%d", int(opts.TTL.Seconds())),
		},
	}

	_, err := c.Logical().Write(metaPath, body)
	if err != nil {
		return ExpireResult{
			Path:    path,
			TTL:     opts.TTL,
			Success: false,
			Err:     fmt.Errorf("failed to write expiry metadata: %w", err),
		}
	}

	return ExpireResult{
		Path:      path,
		TTL:       opts.TTL,
		ExpiresAt: expiresAt,
		Success:   true,
	}
}

// IsExpired returns true if the given time is in the past.
func IsExpired(expiresAt time.Time) bool {
	return !expiresAt.IsZero() && time.Now().UTC().After(expiresAt)
}
