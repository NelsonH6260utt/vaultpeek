package vault

import "fmt"

// LintResult holds the outcome of linting a secret's keys.
type LintResult struct {
	Path     string
	Warnings []string
	Success  bool
	Err      error
}

// LintRule defines a function that inspects a key-value pair and returns a warning or empty string.
type LintRule func(key, value string) string

// DefaultLintRules returns the standard set of lint rules.
func DefaultLintRules() []LintRule {
	return []LintRule{
		func(key, value string) string {
			if value == "" {
				return fmt.Sprintf("key %q has an empty value", key)
			}
			return ""
		},
		func(key, _ string) string {
			if key == "" {
				return "empty key name detected"
			}
			return ""
		},
		func(key, value string) string {
			if len(value) > 4096 {
				return fmt.Sprintf("key %q value exceeds 4096 characters", key)
			}
			return ""
		},
	}
}

// LintSecret runs lint rules against a secret's data map.
func LintSecret(path string, data map[string]interface{}, rules []LintRule) LintResult {
	if data == nil {
		return LintResult{Path: path, Success: false, Err: fmt.Errorf("nil data for path %q", path)}
	}

	var warnings []string
	for k, v := range data {
		val, _ := v.(string)
		for _, rule := range rules {
			if w := rule(k, val); w != "" {
				warnings = append(warnings, w)
			}
		}
	}

	return LintResult{
		Path:     path,
		Warnings: warnings,
		Success:  true,
	}
}
