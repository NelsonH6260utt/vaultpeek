package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultpeek/internal/vault"
)

// PrintLockResult writes a human-readable lock/unlock result to stdout.
func PrintLockResult(r vault.LockResult) {
	printLockResultTo(os.Stdout, r)
}

func printLockResultTo(w io.Writer, r vault.LockResult) {
	if !r.Success {
		if r.Error != nil {
			fmt.Fprintf(w, "[error] lock operation failed for %q: %v\n", r.Path, r.Error)
		} else {
			fmt.Fprintf(w, "[error] lock operation failed for %q\n", r.Path)
		}
		return
	}
	if r.Locked {
		fmt.Fprintf(w, "[locked]   %s\n", r.Path)
		fmt.Fprintf(w, "           secret is now protected from writes\n")
	} else {
		fmt.Fprintf(w, "[unlocked] %s\n", r.Path)
		fmt.Fprintf(w, "           secret lock has been removed\n")
	}
}
