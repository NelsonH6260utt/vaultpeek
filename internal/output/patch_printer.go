package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/vaultpeek/internal/vault"
)

// PrintPatchResult writes a human-readable summary of a PatchResult to stdout.
func PrintPatchResult(result vault.PatchResult) {
	printPatchResultTo(os.Stdout, result)
}

func printPatchResultTo(w io.Writer, result vault.PatchResult) {
	if !result.Success {
		fmt.Fprintf(w, "✗ patch failed for %q\n", result.Path)
		if result.Error != nil {
			fmt.Fprintf(w, "  error: %s\n", result.Error.Error())
		}
		return
	}

	fmt.Fprintf(w, "✔ patched %q\n", result.Path)

	if len(result.UpdatedKeys) == 0 {
		fmt.Fprintln(w, "  no keys updated")
		return
	}

	sorted := make([]string, len(result.UpdatedKeys))
	copy(sorted, result.UpdatedKeys)
	sort.Strings(sorted)

	fmt.Fprintf(w, "  updated %d key(s):\n", len(sorted))
	for _, k := range sorted {
		fmt.Fprintf(w, "    ~ %s\n", k)
	}
}
