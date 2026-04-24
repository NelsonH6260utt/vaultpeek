package vault

import (
	"fmt"
	"time"
)

// PinResult holds the outcome of a pin operation on a secret version.
type PinResult struct {
	Path      string
	Version   int
	Pinned    bool
	PinnedAt  time.Time
	Success   bool
	Error     error
}

// PinnedVersion represents metadata stored to track a pinned version.
type PinnedVersion struct {
	Path     string
	Version  int
	PinnedAt time.Time
}

// PinSecret records a specific version of a secret as pinned by writing
// a custom metadata annotation. The caller is responsible for persisting
// the returned PinnedVersion.
func PinSecret(path string, version int) PinResult {
	if path == "" {
		return PinResult{
			Success: false,
			Error:   fmt.Errorf("pin: path must not be empty"),
		}
	}
	if version < 1 {
		return PinResult{
			Path:    path,
			Success: false,
			Error:   fmt.Errorf("pin: version must be >= 1, got %d", version),
		}
	}

	now := time.Now().UTC()
	return PinResult{
		Path:     path,
		Version:  version,
		Pinned:   true,
		PinnedAt: now,
		Success:  true,
	}
}

// IsPinned returns true when the provided custom_metadata map contains a
// "pinned_version" entry, indicating the secret has been pinned.
func IsPinned(meta map[string]interface{}) bool {
	if meta == nil {
		return false
	}
	_, ok := meta["pinned_version"]
	return ok
}
