package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

// VersionDiffResult holds the comparison between two versions of a secret.
type VersionDiffResult struct {
	Path     string
	VersionA int
	VersionB int
	Added    map[string]string
	Removed  map[string]string
	Changed  map[string][2]string
	Success  bool
	Err      error
}

// DiffVersions compares two versions of a KVv2 secret at the given path.
func DiffVersions(client *api.Client, mount, path string, versionA, versionB int) VersionDiffResult {
	result := VersionDiffResult{
		Path:     path,
		VersionA: versionA,
		VersionB: versionB,
		Added:    make(map[string]string),
		Removed:  make(map[string]string),
		Changed:  make(map[string][2]string),
	}

	fetchVersion := func(version int) (map[string]string, error) {
		dataPath := KVv2DataPath(mount, path)
		secret, err := client.Logical().ReadWithData(dataPath, map[string][]string{
			"version": {fmt.Sprintf("%d", version)},
		})
		if err != nil {
			return nil, err
		}
		if secret == nil || secret.Data == nil {
			return nil, fmt.Errorf("version %d not found at %s", version, path)
		}
		raw, ok := secret.Data["data"]
		if !ok {
			return nil, fmt.Errorf("no data field in version %d", version)
		}
		m, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected data format in version %d", version)
		}
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[k] = fmt.Sprintf("%v", v)
		}
		return out, nil
	}

	dataA, err := fetchVersion(versionA)
	if err != nil {
		result.Err = err
		return result
	}

	dataB, err := fetchVersion(versionB)
	if err != nil {
		result.Err = err
		return result
	}

	for k, va := range dataA {
		if vb, exists := dataB[k]; exists {
			if va != vb {
				result.Changed[k] = [2]string{va, vb}
			}
		} else {
			result.Removed[k] = va
		}
	}
	for k, vb := range dataB {
		if _, exists := dataA[k]; !exists {
			result.Added[k] = vb
		}
	}

	result.Success = true
	return result
}
