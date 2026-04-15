package output

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultpeek/internal/vault"
)

// PrintRenameResult writes the result of a rename operation to stdout.
func PrintRenameResult(result vault.RenameResult) {
	printRenameResultTo(os.Stdout, result)
}

func printRenameResultTo(w io.Writer, result vault.RenameResult) {
	if !result.Success {
		fmt.Fprintf(w, "✗ Rename failed: %s\n", result.Error)
		return
	}

	fmt.Fprintf(w, "✓ Renamed secret\n")
	fmt.Fprintf(w, "  From : %s\n", result.SourcePath)
	fmt.Fprintf(w, "  To   : %s\n", result.DestPath)

	if len(result.Keys) == 0 {
		fmt.Fprintln(w, "  Keys : (none)")
		return
	}

	fmt.Fprintf(w, "  Keys : %d transferred\n", len(result.Keys))
	for _, k := range result.Keys {
		fmt.Fprintf(w, "    - %s\n", k)
	}
}
