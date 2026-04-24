package vault

import (
	"errors"
	"sort"
	"time"
)

// AnnotateResult holds the outcome of an annotation operation on a secret path.
type AnnotateResult struct {
	Path        string
	Annotations map[string]string
	UpdatedAt   time.Time
	Success     bool
	Err         error
}

// SortedAnnotationKeys returns the annotation keys in sorted order.
func SortedAnnotationKeys(annotations map[string]string) []string {
	keys := make([]string, 0, len(annotations))
	for k := range annotations {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// AnnotateSecret writes free-form annotation key/value pairs into the
// custom_metadata of a KV v2 secret. Annotations are merged with any
// existing custom_metadata entries.
func AnnotateSecret(client *Client, path string, annotations map[string]string) AnnotateResult {
	if path == "" {
		return AnnotateResult{
			Path:    path,
			Success: false,
			Err:     errors.New("path must not be empty"),
		}
	}

	if len(annotations) == 0 {
		return AnnotateResult{
			Path:    path,
			Success: false,
			Err:     errors.New("annotations must not be empty"),
		}
	}

	metaPath := KVv2MetaPath(path)

	// Read existing custom_metadata so we can merge.
	existing := map[string]interface{}{}
	secret, err := client.Logical().Read(metaPath)
	if err == nil && secret != nil {
		if cm, ok := secret.Data["custom_metadata"]; ok {
			if cmMap, ok := cm.(map[string]interface{}); ok {
				for k, v := range cmMap {
					existing[k] = v
				}
			}
		}
	}

	for k, v := range annotations {
		existing[k] = v
	}

	_, writeErr := client.Logical().Write(metaPath, map[string]interface{}{
		"custom_metadata": existing,
	})
	if writeErr != nil {
		return AnnotateResult{
			Path:    path,
			Success: false,
			Err:     writeErr,
		}
	}

	merged := make(map[string]string, len(existing))
	for k, v := range existing {
		if sv, ok := v.(string); ok {
			merged[k] = sv
		}
	}

	return AnnotateResult{
		Path:        path,
		Annotations: merged,
		UpdatedAt:   time.Now().UTC(),
		Success:     true,
	}
}
