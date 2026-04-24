package vault

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

// SignResult holds the outcome of a secret signing operation.
type SignResult struct {
	Path      string
	Signature string
	Algorithm string
	SignedAt  time.Time
	KeyCount  int
	Success   bool
	Err       error
}

// SignSecret computes an HMAC-SHA256 signature over the canonical form of a
// secret's key-value data using the provided signingKey. The canonical form
// is built by sorting keys lexicographically and joining them as
// "key=value\n" lines, ensuring deterministic output.
func SignSecret(path string, data map[string]interface{}, signingKey string) SignResult {
	if data == nil {
		return SignResult{
			Path:    path,
			Success: false,
			Err:     fmt.Errorf("sign: nil data for path %q", path),
		}
	}

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%v\n", k, data[k])
	}

	mac := hmac.New(sha256.New, []byte(signingKey))
	mac.Write([]byte(sb.String()))
	sig := hex.EncodeToString(mac.Sum(nil))

	return SignResult{
		Path:      path,
		Signature: sig,
		Algorithm: "hmac-sha256",
		SignedAt:  time.Now().UTC(),
		KeyCount:  len(keys),
		Success:   true,
	}
}

// VerifySecret recomputes the signature for data and compares it to the
// expected signature. Returns true only when the signatures match.
func VerifySecret(data map[string]interface{}, signingKey, expected string) bool {
	if data == nil || signingKey == "" || expected == "" {
		return false
	}
	result := SignSecret("", data, signingKey)
	if !result.Success {
		return false
	}
	return hmac.Equal([]byte(result.Signature), []byte(expected))
}
