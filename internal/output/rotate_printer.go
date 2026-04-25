package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/your-org/vaultpeek/internal/vault"
)

// PrintRotateResult writes a human-readable summary of a RotateResult to stdout.
func PrintRotateResult(r vault.RotateResult) {
	printRotateResultTo(os.Stdout, r)
}

func printRotateResultTo(w io.Writer, r vault.RotateResult) {
	if !r.Success {
		fmt.Fprintf(w, "[error] rotate failed for %s", r.Path)
		if r.Err != nil {
			fmt.Fprintf(w, ": %s", r.Err.Error())
		}
		fmt.Fprintln(w)
		return
	}

	fmt.Fprintf(w, "[rotated] %s\n", r.Path)
	fmt.Fprintf(w, "  version : %d -> %d\n", r.OldVersion, r.NewVersion)
	fmt.Fprintf(w, "  rotated : %s\n", r.RotatedAt.Format("2006-01-02 15:04:05 UTC"))

	if len(r.Keys) == 0 {
		fmt.Fprintln(w, "  keys    : (none)")
		return
	}

	sorted := make([]string, len(r.Keys))
	copy(sorted, r.Keys)
	sort.Strings(sorted)

	fmt.Fprintf(w, "  keys    : %d rotated\n", len(sorted))
	for _, k := range sorted {
		fmt.Fprintf(w, "    - %s\n", k)
	}
}
