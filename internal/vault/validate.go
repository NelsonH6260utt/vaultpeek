package vault

import (
	"fmt"
	"strings"
)

// ValidationResult holds the outcome of a secret validation check.
type ValidationResult struct {
	Path    string
	Valid   bool
	Errors  []string
	Checked int
}

// ValidationRule is a function that validates a key-value pair.
type ValidationRule func(key, value string) error

// DefaultValidationRules returns a standard set of validation rules.
func DefaultValidationRules() []ValidationRule {
	return []ValidationRule{
		func(key, value string) error {
			if strings.TrimSpace(key) == "" {
				return fmt.Errorf("empty key found")
			}
			return nil
		},
		func(key, value string) error {
			if strings.TrimSpace(value) == "" {
				return fmt.Errorf("key %q has empty or whitespace-only value", key)
			}
			return nil
		},
		func(key, value string) error {
			if strings.Contains(key, " ") {
				return fmt.Errorf("key %q contains spaces", key)
			}
			return nil
		},
	}
}

// ValidateSecret runs validation rules against secret data at the given path.
func ValidateSecret(path string, data map[string]interface{}, rules []ValidationRule) ValidationResult {
	result := ValidationResult{
		Path:  path,
		Valid: true,
	}

	if data == nil {
		result.Valid = false
		result.Errors = append(result.Errors, "secret data is nil")
		return result
	}

	for k, v := range data {
		val, _ := v.(string)
		for _, rule := range rules {
			if err := rule(k, val); err != nil {
				result.Errors = append(result.Errors, err.Error())
				result.Valid = false
			}
		}
		result.Checked++
	}

	return result
}
