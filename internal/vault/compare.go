package vault

import (
	"context"
	"fmt"
)

// SecretMap is a map of key-value pairs from a Vault secret.
type SecretMap map[string]interface{}

// FetchSecretData retrieves the key-value data for a KVv2 secret at the given
// mount and path using the provided client.
func FetchSecretData(ctx context.Context, c *Client, mount, path string) (SecretMap, error) {
	dataPath := KVv2DataPath(mount, path)
	secret, err := c.Logical().ReadWithContext(ctx, dataPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret %q: %w", dataPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("secret not found at %q", dataPath)
	}

	data, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("secret %q has no data field", dataPath)
	}

	kvData, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("secret %q data field is not a map", dataPath)
	}

	return SecretMap(kvData), nil
}

// FetchAndCompare fetches the secret at path from two different clients
// (representing two environments) and returns both SecretMaps for diffing.
func FetchAndCompare(
	ctx context.Context,
	leftClient *Client, leftMount string,
	rightClient *Client, rightMount string,
	path string,
) (SecretMap, SecretMap, error) {
	left, err := FetchSecretData(ctx, leftClient, leftMount, path)
	if err != nil {
		return nil, nil, fmt.Errorf("left env: %w", err)
	}

	right, err := FetchSecretData(ctx, rightClient, rightMount, path)
	if err != nil {
		return nil, nil, fmt.Errorf("right env: %w", err)
	}

	return left, right, nil
}

// DiffKeys returns three slices: keys only in left, keys only in right, and
// keys present in both maps.
func DiffKeys(left, right SecretMap) (onlyLeft, onlyRight, shared []string) {
	for k := range left {
		if _, ok := right[k]; ok {
			shared = append(shared, k)
		} else {
			onlyLeft = append(onlyLeft, k)
		}
	}
	for k := range right {
		if _, ok := left[k]; !ok {
			onlyRight = append(onlyRight, k)
		}
	}
	return onlyLeft, onlyRight, shared
}
