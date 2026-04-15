package vault

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// ExportResult holds the outcome of an export operation.
type ExportResult struct {
	Path    string
	Data    map[string]interface{}
	Success bool
	Err     error
}

// ExportSecret reads a KVv2 secret at the given mount and path,
// returning a structured result suitable for serialization.
func ExportSecret(client *api.Client, mount, path string) ExportResult {
	dataPath := KVv2DataPath(mount, path)
	secret, err := client.Logical().Read(dataPath)
	if err != nil {
		return ExportResult{Path: path, Success: false, Err: err}
	}
	if secret == nil || secret.Data == nil {
		return ExportResult{Path: path, Success: false, Err: fmt.Errorf("no secret found at %s", path)}
	}

	kv, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return ExportResult{Path: path, Success: false, Err: fmt.Errorf("unexpected data format at %s", path)}
	}

	return ExportResult{Path: path, Data: kv, Success: true}
}

// MarshalExport serializes an ExportResult's Data to indented JSON bytes.
func MarshalExport(result ExportResult) ([]byte, error) {
	if !result.Success {
		return nil, result.Err
	}
	return json.MarshalIndent(result.Data, "", "  ")
}
