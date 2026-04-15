package diff

import (
	"fmt"
	"sort"
	"strings"
)

// Result holds the comparison between two sets of secret key-value pairs.
type Result struct {
	OnlyInLeft  map[string]interface{}
	OnlyInRight map[string]interface{}
	Different   map[string][2]interface{}
	Identical   map[string]interface{}
}

// Compare compares two secret data maps and returns a Result.
func Compare(left, right map[string]interface{}) Result {
	result := Result{
		OnlyInLeft:  make(map[string]interface{}),
		OnlyInRight: make(map[string]interface{}),
		Different:   make(map[string][2]interface{}),
		Identical:   make(map[string]interface{}),
	}

	for k, lv := range left {
		if rv, ok := right[k]; ok {
			if fmt.Sprintf("%v", lv) == fmt.Sprintf("%v", rv) {
				result.Identical[k] = lv
			} else {
				result.Different[k] = [2]interface{}{lv, rv}
			}
		} else {
			result.OnlyInLeft[k] = lv
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			result.OnlyInRight[k] = rv
		}
	}

	return result
}

// HasDifferences returns true if there are any differences between the two secrets.
func (r Result) HasDifferences() bool {
	return len(r.OnlyInLeft) > 0 || len(r.OnlyInRight) > 0 || len(r.Different) > 0
}

// Summary returns a human-readable summary of the diff result.
func (r Result) Summary(leftLabel, rightLabel string) string {
	var sb strings.Builder

	keys := func(m map[string]interface{}) []string {
		out := make([]string, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		sort.Strings(out)
		return out
	}

	for _, k := range keys(r.OnlyInLeft) {
		sb.WriteString(fmt.Sprintf("< [%s only] %s\n", leftLabel, k))
	}
	for _, k := range keys(r.OnlyInRight) {
		sb.WriteString(fmt.Sprintf("> [%s only] %s\n", rightLabel, k))
	}
	for _, k := range keys(r.Different) {
		pair := r.Different[k]
		sb.WriteString(fmt.Sprintf("~ [changed] %s: %v -> %v\n", k, pair[0], pair[1]))
	}

	if !r.HasDifferences() {
		sb.WriteString("No differences found.\n")
	}

	return sb.String()
}
