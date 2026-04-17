package output

import (
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpeek/internal/vault"
)

// PrintLintResult writes a lint result to stdout.
func PrintLintResult(r vault.LintResult) {
	printLintResultTo(os.Stdout, r)
}

func printLintResultTo(w io.Writer, r vault.LintResult) {
	if !r.Success {
		fmt.Fprintf(w, "[lint] error: %v\n", r.Err)
		return
	}

	if len(r.Warnings) == 0 {
		fmt.Fprintf(w, "[lint] %s — no issues found\n", r.Path)
		return
	}

	fmt.Fprintf(w, "[lint] %s — %d warning(s):\n", r.Path, len(r.Warnings))
	for _, w2 := range r.Warnings {
		fmt.Fprintf(w, "  ⚠  %s\n", w2)
	}
}
