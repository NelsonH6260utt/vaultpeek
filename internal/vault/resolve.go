package vault

import (
	"fmt"
	"strings"
)

// ResolveResult holds the outcome of resolving a secret path alias or reference.
type ResolveResult struct {
	Path     string
	Resolved string
	Aliases  map[string]string
	Success  bool
	Err      error
}

// ResolveOptions configures path resolution behaviour.
type ResolveOptions struct {
	// Aliases maps short names to full Vault paths.
	Aliases map[string]string
	// Mount is the KV mount prefix to prepend when no alias matches.
	Mount string
}

// ResolveSecret resolves a path or alias to a canonical Vault secret path.
// If the input matches a key in opts.Aliases, the mapped path is returned.
// Otherwise the path is normalised against opts.Mount.
func ResolveSecret(path string, opts ResolveOptions) ResolveResult {
	if path == "" {
		return ResolveResult{
			Path:    path,
			Success: false,
			Err:     fmt.Errorf("resolve: path must not be empty"),
		}
	}

	if opts.Aliases != nil {
		if target, ok := opts.Aliases[path]; ok {
			return ResolveResult{
				Path:     path,
				Resolved: target,
				Aliases:  opts.Aliases,
				Success:  true,
			}
		}
	}

	mount := strings.TrimRight(opts.Mount, "/")
	resolved := path
	if mount != "" && !strings.HasPrefix(path, mount+"/") {
		resolved = mount + "/" + strings.TrimLeft(path, "/")
	}

	return ResolveResult{
		Path:     path,
		Resolved: resolved,
		Aliases:  opts.Aliases,
		Success:  true,
	}
}
