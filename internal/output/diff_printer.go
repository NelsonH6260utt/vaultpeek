package output

import (
	"fmt"
	"sort"

	"github.com/user/vaultpeek/internal/diff"
)

// PrintDiff renders a diff.Result to the formatter's writer.
// leftLabel and rightLabel name the two environments being compared.
func (f *Formatter) PrintDiff(result diff.Result, leftLabel, rightLabel string) {
	fmt.Fprintf(f.w, "Comparing: %s  <>  %s\n", leftLabel, rightLabel)
	fmt.Fprintf(f.w, "%s\n", repeatDash(44))

	keys := collectKeys(result)
	for _, k := range keys {
		if v, ok := result.OnlyInLeft[k]; ok {
			f.PrintDiffLine("-", k, v)
			continue
		}
		if v, ok := result.OnlyInRight[k]; ok {
			f.PrintDiffLine("+", k, v)
			continue
		}
		if ch, ok := result.Different[k]; ok {
			f.PrintDiffLine("-", k, ch.Left)
			f.PrintDiffLine("+", k, ch.Right)
			continue
		}
		if v, ok := result.Identical[k]; ok {
			f.PrintDiffLine(" ", k, v)
		}
	}

	fmt.Fprintf(f.w, "%s\n", repeatDash(44))
	fmt.Fprintf(f.w, "Summary: %d identical, %d only-left, %d only-right, %d changed\n",
		len(result.Identical),
		len(result.OnlyInLeft),
		len(result.OnlyInRight),
		len(result.Different),
	)
}

func repeatDash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = '-'
	}
	return string(b)
}

func collectKeys(r diff.Result) []string {
	seen := make(map[string]struct{})
	for k := range r.Identical {
		seen[k] = struct{}{}
	}
	for k := range r.OnlyInLeft {
		seen[k] = struct{}{}
	}
	for k := range r.OnlyInRight {
		seen[k] = struct{}{}
	}
	for k := range r.Different {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
