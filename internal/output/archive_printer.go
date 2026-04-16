package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/vaultpeek/internal/vault"
)

// PrintArchiveResult writes the result of an archive operation to stdout.
func PrintArchiveResult(result vault.ArchiveResult, entry vault.ArchiveEntry) {
	printArchiveResultTo(os.Stdout, result, entry)
}

func printArchiveResultTo(w io.Writer, result vault.ArchiveResult, entry vault.ArchiveEntry) {
	if !result.Success {
		fmt.Fprintf(w, "✗ Archive failed for %q", result.Path)
		if result.Error != nil {
			fmt.Fprintf(w, ": %s", result.Error)
		}
		fmt.Fprintln(w)
		return
	}

	fmt.Fprintf(w, "✔ Archived %q at %s\n", entry.Path, entry.ArchivedAt.Format("2006-01-02T15:04:05Z"))

	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "  Keys archived (%d):\n", len(keys))
	for _, k := range keys {
		fmt.Fprintf(w, "    - %s\n", k)
	}
}
