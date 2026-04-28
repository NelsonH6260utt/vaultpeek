package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultpeek/internal/vault"
)

// PrintResolveResult writes the outcome of a path resolution to stdout.
func PrintResolveResult(r vault.ResolveResult) {
	printResolveResultTo(os.Stdout, r)
}

func printResolveResultTo(w io.Writer, r vault.ResolveResult) {
	if !r.Success {
		fmt.Fprintf(w, "[error] resolve failed: %v\n", r.Err)
		return
	}

	if r.Path == r.Resolved {
		fmt.Fprintf(w, "[resolve] %s (no change)\n", r.Resolved)
	} else {
		fmt.Fprintf(w, "[resolve] %s => %s\n", r.Path, r.Resolved)
	}

	if len(r.Aliases) == 0 {
		return
	}

	keys := make([]string, 0, len(r.Aliases))
	for k := range r.Aliases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(w, "  known aliases:")
	for _, k := range keys {
		fmt.Fprintf(w, "    %-20s => %s\n", k, r.Aliases[k])
	}
}
