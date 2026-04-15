package output

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpeek/internal/vault"
)

// PrintRollbackResult writes a human-readable rollback outcome to stdout.
func PrintRollbackResult(result vault.RollbackResult) {
	printRollbackResultTo(os.Stdout, result)
}

func printRollbackResultTo(w io.Writer, result vault.RollbackResult) {
	if !result.Success {
		fmt.Fprintf(w, "✗ Rollback failed for %q\n", result.Path)
		if result.Error != nil {
			fmt.Fprintf(w, "  Error: %s\n", result.Error.Error())
		}
		return
	}

	fmt.Fprintf(w, "✔ Rolled back %q\n", result.Path)
	fmt.Fprintf(w, "  From version : %d\n", result.FromVersion)
	fmt.Fprintf(w, "  To version   : %d\n", result.ToVersion)
	fmt.Fprintf(w, "  New version  : %d (written as latest)\n", result.FromVersion+1)
}
