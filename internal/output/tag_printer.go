package output

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpeek/internal/vault"
)

// PrintTagResult writes a human-readable tag result to stdout.
func PrintTagResult(result vault.TagResult) {
	printTagResultTo(os.Stdout, result)
}

func printTagResultTo(w io.Writer, result vault.TagResult) {
	if !result.Success {
		fmt.Fprintf(w, "[error] failed to tag %s", result.Path)
		if result.Err != nil {
			fmt.Fprintf(w, ": %s", result.Err.Error())
		}
		fmt.Fprintln(w)
		return
	}

	fmt.Fprintf(w, "Tagged: %s\n", result.Path)
	if len(result.Tags) == 0 {
		fmt.Fprintln(w, "  (no tags)")
		return
	}
	for _, k := range vault.SortedTagKeys(result.Tags) {
		fmt.Fprintf(w, "  %s = %s\n", k, result.Tags[k])
	}
}
