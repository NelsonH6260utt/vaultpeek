package vault

import (
	"strings"
)

// SanitizeResult holds the outcome of a sanitize operation.
type SanitizeResult struct {
	Path    string
	Success bool
	Error   error
	Changed map[string]string // key -> old value (redacted)
	Keys    []string
}

// SanitizeRule defines a function that determines whether a key's value should be sanitized.
type SanitizeRule func(key, value string) bool

// DefaultSanitizeRules returns rules that flag common sensitive key patterns.
func DefaultSanitizeRules() []SanitizeRule {
	sensitivePatterns := []string{"password", "secret", "token", "apikey", "api_key", "private_key"}
	return []SanitizeRule{
		func(key, _ string) bool {
			lower := strings.ToLower(key)
			for _, p := range sensitivePatterns {
				if strings.Contains(lower, p) {
					return true
				}
			}
			return false
		},
	}
}

// SanitizeSecret redacts values in data that match any of the provided rules.
// It returns a SanitizeResult describing what was changed without modifying the original map.
func SanitizeSecret(path string, data map[string]interface{}, rules []SanitizeRule) SanitizeResult {
	if data == nil {
		return SanitizeResult{Path: path, Success: true}
	}

	changed := make(map[string]string)
	keys := make([]string, 0, len(data))

	for k, v := range data {
		keys = append(keys, k)
		val, _ := v.(string)
		for _, rule := range rules {
			if rule(k, val) {
				changed[k] = "[REDACTED]"
				break
			}
		}
	}

	return SanitizeResult{
		Path:    path,
		Success: true,
		Changed: changed,
		Keys:    keys,
	}
}
