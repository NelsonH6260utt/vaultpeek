// Package diff provides utilities for comparing two sets of Vault secret
// key-value pairs retrieved from different environments or paths.
//
// Use Compare to produce a Result that categorizes keys as identical,
// changed, only present on the left side, or only present on the right side.
//
// Example usage:
//
//	result := diff.Compare(stagingSecrets, prodSecrets)
//	if result.HasDifferences() {
//		fmt.Print(result.Summary("staging", "prod"))
//	}
package diff
