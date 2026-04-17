package vault

import "fmt"

// MergeResult holds the outcome of merging two secrets.
type MergeResult struct {
	Path    string
	Keys    []string
	Success bool
	Err     error
}

// MergeStrategy controls how key conflicts are resolved.
type MergeStrategy int

const (
	// MergeStrategyKeepDest keeps the destination value on conflict.
	MergeStrategyKeepDest MergeStrategy = iota
	// MergeStrategyOverwrite overwrites destination with source value.
	MergeStrategyOverwrite
)

// MergeSecret merges keys from src into dst at destPath using the given strategy.
// The merged result is written back via the provided write function.
func MergeSecret(
	srcData map[string]interface{},
	dstData map[string]interface{},
	destPath string,
	strategy MergeStrategy,
	write func(path string, data map[string]interface{}) error,
) MergeResult {
	if srcData == nil {
		return MergeResult{Path: destPath, Success: false, Err: fmt.Errorf("source data is nil")}
	}
	if dstData == nil {
		dstData = make(map[string]interface{})
	}

	merged := make(map[string]interface{}, len(dstData))
	for k, v := range dstData {
		merged[k] = v
	}

	for k, v := range srcData {
		if _, exists := merged[k]; exists {
			if strategy == MergeStrategyOverwrite {
				merged[k] = v
			}
			// MergeStrategyKeepDest: do nothing
		} else {
			merged[k] = v
		}
	}

	if err := write(destPath, merged); err != nil {
		return MergeResult{Path: destPath, Success: false, Err: err}
	}

	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}

	return MergeResult{Path: destPath, Keys: keys, Success: true}
}
