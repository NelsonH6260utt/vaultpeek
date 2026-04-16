package vault

import (
	"fmt"
	"time"
)

// LockResult holds the outcome of a lock or unlock operation.
type LockResult struct {
	Path    string
	Locked  bool
	Success bool
	Error   error
}

// LockMeta holds lock metadata stored in Vault secret metadata custom fields.
type LockMeta struct {
	LockedBy  string
	LockedAt  time.Time
	Reason    string
}

// LockSecret marks a secret path as locked by writing lock metadata.
func LockSecret(c *Client, mount, path, lockedBy, reason string) LockResult {
	if c == nil {
		return LockResult{Path: path, Success: false, Error: fmt.Errorf("nil client")}
	}
	tags := map[string]string{
		"locked":    "true",
		"locked_by": lockedBy,
		"locked_at": time.Now().UTC().Format(time.RFC3339),
		"lock_reason": reason,
	}
	_, err := c.Logical().Write(
		KVv2MetaPath(mount, path),
		map[string]interface{}{"custom_metadata": tags},
	)
	if err != nil {
		return LockResult{Path: path, Locked: false, Success: false, Error: err}
	}
	return LockResult{Path: path, Locked: true, Success: true}
}

// UnlockSecret removes the lock metadata from a secret path.
func UnlockSecret(c *Client, mount, path string) LockResult {
	if c == nil {
		return LockResult{Path: path, Success: false, Error: fmt.Errorf("nil client")}
	}
	tags := map[string]string{
		"locked":      "",
		"locked_by":   "",
		"locked_at":   "",
		"lock_reason": "",
	}
	_, err := c.Logical().Write(
		KVv2MetaPath(mount, path),
		map[string]interface{}{"custom_metadata": tags},
	)
	if err != nil {
		return LockResult{Path: path, Locked: true, Success: false, Error: err}
	}
	return LockResult{Path: path, Locked: false, Success: true}
}

// IsLocked returns true if the secret's custom metadata marks it as locked.
func IsLocked(meta map[string]interface{}) bool {
	cm, ok := meta["custom_metadata"].(map[string]interface{})
	if !ok {
		return false
	}
	v, _ := cm["locked"].(string)
	return v == "true"
}
