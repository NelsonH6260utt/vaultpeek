package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultpeek/internal/vault"
)

// PrintValidateResult prints the result of a secret validation to stdout.
func PrintValidateResult(result vault.ValidationResult) {
	printValidateResultTo(os.Stdout, result)
}

func printValidateResultTo(w io.Writer, result vault.ValidationResult) {
	fmt.Fprintf(w, "Path: %s\n", result.Path)
	fmt.Fprintf(w, "Keys checked: %d\n", result.Checked)

	if result.Valid {
		fmt.Fprintln(w, "Status: ✓ valid")
		return
	}

	fmt.Fprintln(w, "Status: ✗ invalid")
	fmt.Fprintf(w, "Errors (%d):\n", len(result.Errors))
	for _, e := range result.Errors {
		fmt.Fprintf(w, "  - %s\n", e)
	}
}
