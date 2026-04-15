package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/yourusername/vaultpeek/internal/diff"
	"github.com/yourusername/vaultpeek/internal/vault"
)

// PrintCompareSummary writes a human-readable diff summary between two
// SecretMaps to w, labelled with leftEnv and rightEnv.
func PrintCompareSummary(w io.Writer, path, leftEnv, rightEnv string, left, right vault.SecretMap) {
	fmt.Fprintf(w, "\nComparing secret: %s\n", path)
	fmt.Fprintf(w, "  ← %s  |  → %s\n", leftEnv, rightEnv)
	fmt.Fprintln(w, repeatDash(48))

	leftRaw := toRawMap(left)
	rightRaw := toRawMap(right)

	results := diff.Compare(leftRaw, rightRaw)

	keys := collectResultKeys(results)
	sort.Strings(keys)

	if len(keys) == 0 {
		fmt.Fprintln(w, "  ✓ No differences found.")
		return
	}

	for _, k := range keys {
		r := results[k]
		switch {
		case r.OnlyInLeft:
			fmt.Fprintf(w, "  - [%s only in %s] %s = %v\n", leftEnv, leftEnv, k, r.LeftValue)
		case r.OnlyInRight:
			fmt.Fprintf(w, "  + [%s only in %s] %s = %v\n", rightEnv, rightEnv, k, r.RightValue)
		case r.Different:
			fmt.Fprintf(w, "  ~ [changed] %s\n", k)
			fmt.Fprintf(w, "      ← %v\n", r.LeftValue)
			fmt.Fprintf(w, "      → %v\n", r.RightValue)
		}
	}

	fmt.Fprintln(w, repeatDash(48))
	fmt.Fprintln(w, diff.Summary(results))
}

func toRawMap(sm vault.SecretMap) map[string]interface{} {
	out := make(map[string]interface{}, len(sm))
	for k, v := range sm {
		out[k] = v
	}
	return out
}

func collectResultKeys(results map[string]diff.Result) []string {
	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	return keys
}
