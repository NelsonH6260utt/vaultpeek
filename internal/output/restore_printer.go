package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/vaultpeek/internal/vault"
)

// PrintRestoreResult prints the outcome of a secret restore to stdout.
func PrintRestoreResult(r vault.RestoreResult) {
	printRestoreResultTo(os.Stdout, r)
}

func printRestoreResultTo(w io.Writer, r vault.RestoreResult) {
	if !r.Success {
		if r.Err != nil {
			fmt.Fprintf(w, "[error] restore failed for %s: %v\n", r.Path, r.Err)
		} else {
			fmt.Fprintf(w, "[error] restore failed for %s\n", r.Path)
		}
		return
	}

	fmt.Fprintf(w, "[restored] %s\n", r.Path)
	fmt.Fprintf(w, "  from version : %d\n", r.FromVersion)
	fmt.Fprintf(w, "  new version  : %d\n", r.NewVersion)
	fmt.Fprintf(w, "  restored at  : %s\n", r.RestoredAt.Format("2006-01-02 15:04:05 UTC"))

	if len(r.Keys) == 0 {
		fmt.Fprintf(w, "  keys         : (none)\n")
		return
	}

	sorted := make([]string, len(r.Keys))
	copy(sorted, r.Keys)
	sort.Strings(sorted)

	fmt.Fprintf(w, "  keys         : %d\n", len(sorted))
	for _, k := range sorted {
		fmt.Fprintf(w, "    - %s\n", k)
	}
}

// PrintRestoreResults prints the outcome of multiple secret restores to stdout.
func PrintRestoreResults(results []vault.RestoreResult) {
	for _, r := range results {
		PrintRestoreResult(r)
	}

	successes := 0
	for _, r := range results {
		if r.Success {
			successes++
		}
	}
	fmt.Fprintf(os.Stdout, "\nsummary: %d/%d restored successfully\n", successes, len(results))
}
