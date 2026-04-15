package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultpeek/internal/vault"
)

// PrintProtectResult writes a human-readable summary of a ProtectResult to
// stdout. On failure the error is written to stderr.
func PrintProtectResult(result vault.ProtectResult) {
	printProtectResultTo(os.Stdout, os.Stderr, result)
}

func printProtectResultTo(out, errOut io.Writer, result vault.ProtectResult) {
	if !result.Success {
		fmt.Fprintf(errOut, "error: could not protect secret at %q\n", result.Path)
		if result.Error != nil {
			fmt.Fprintf(errOut, "  reason: %v\n", result.Error)
		}
		return
	}

	fmt.Fprintf(out, "protected  %s/%s\n", result.Mount, result.Path)
	fmt.Fprintf(out, "  The path is now marked with vaultpeek_protected=true in custom metadata.\n")
}
