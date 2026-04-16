package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/yourusername/vaultpeek/internal/vault"
)

// PrintSnapshotResult writes a snapshot summary to stdout.
func PrintSnapshotResult(r vault.SnapshotResult) {
	printSnapshotResultTo(os.Stdout, r)
}

func printSnapshotResultTo(w io.Writer, r vault.SnapshotResult) {
	if !r.Success {
		fmt.Fprintf(w, "[error] snapshot failed: %v\n", r.Err)
		return
	}

	fmt.Fprintf(w, "Snapshot taken at: %s\n", r.TakenAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Paths captured: %d\n", len(r.Entries))
	fmt.Fprintln(w, repeatDash(40))

	for _, entry := range r.Entries {
		fmt.Fprintf(w, "  %s\n", entry.Path)
		keys := make([]string, 0, len(entry.Data))
		for k := range entry.Data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "    %s = %v\n", k, entry.Data[k])
		}
	}
}
