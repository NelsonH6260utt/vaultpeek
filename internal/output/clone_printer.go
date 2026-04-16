package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/vaultpeek/internal/vault"
)

// PrintCloneResult writes a human-readable summary of a CloneResult to stdout.
func PrintCloneResult(r vault.CloneResult) {
	printCloneResultTo(os.Stdout, r)
}

func printCloneResultTo(w io.Writer, r vault.CloneResult) {
	if !r.Success {
		if r.Error != nil {
			fmt.Fprintf(w, "✗ Clone failed: %s\n", r.Error.Error())
		} else {
			fmt.Fprintln(w, "✗ Clone failed: unknown error")
		}
		return
	}

	fmt.Fprintf(w, "✔ Cloned %q → %q\n", r.SourcePath, r.DestPath)

	keys := make([]string, len(r.Keys))
	copy(keys, r.Keys)
	sort.Strings(keys)

	fmt.Fprintf(w, "  Keys copied (%d):\n", len(keys))
	for _, k := range keys {
		fmt.Fprintf(w, "    - %s\n", k)
	}
}
