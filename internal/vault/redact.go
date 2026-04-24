package vault

import (
	"regexp"
	"strings"
)

// RedactRule defines a pattern-based rule for redacting secret values.
type RedactRule struct {
	Name    string
	Pattern *regexp.Regexp
	Mask    string
}

// RedactResult holds the outcome of a redaction operation.
type RedactResult struct {
	Path    string
	Redacted map[string]string
	Count   int
	Success bool
	Err     error
}

// DefaultRedactRules returns a set of built-in redaction rules.
func DefaultRedactRules() []RedactRule {
	return []RedactRule{
		{
			Name:    "password",
			Pattern: regexp.MustCompile(`(?i)(password|passwd|pwd)`),
			Mask:    "***REDACTED***",
		},
		{
			Name:    "token",
			Pattern: regexp.MustCompile(`(?i)(token|secret|apikey|api_key)`),
			Mask:    "***REDACTED***",
		},
		{
			Name:    "private_key",
			Pattern: regexp.MustCompile(`(?i)(private[_-]?key|priv[_-]?key)`),
			Mask:    "***REDACTED***",
		},
	}
}

// RedactSecret applies redaction rules to secret data, masking sensitive values.
// It returns a RedactResult with a copy of the data with sensitive fields masked.
func RedactSecret(path string, data map[string]interface{}, rules []RedactRule) RedactResult {
	if data == nil {
		return RedactResult{
			Path:    path,
			Redacted: map[string]string{},
			Success: true,
		}
	}

	redacted := make(map[string]string, len(data))
	count := 0

	for k, v := range data {
		str := ""
		if v != nil {
			str = strings.TrimSpace(fmt.Sprintf("%v", v))
		}

		masked := false
		for _, rule := range rules {
			if rule.Pattern.MatchString(k) {
				redacted[k] = rule.Mask
				count++
				masked = true
				break
			}
		}

		if !masked {
			redacted[k] = str
		}
	}

	return RedactResult{
		Path:    path,
		Redacted: redacted,
		Count:   count,
		Success: true,
	}
}
