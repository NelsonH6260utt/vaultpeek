package output

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/your-org/vaultpeek/internal/vault"
)

// PrintCopyResult writes a human-readable summary of a secret copy operation
// to the provided writer.
func PrintCopyResult(w io.Writer, result vault.CopyResult) {
	if result.Err != nil {
		fmt.Fprintf(w, "[error] copy failed: %v\n", result.Err)
		return
	}

	if !result.Success {
		fmt.Fprintf(w, "[warn] copy did not complete successfully\n")
		return
	}

	fmt.Fprintf(w, "[ok] copied secret\n")
	fmt.Fprintf(w, "     from : %s\n", result.SourcePath)
	fmt.Fprintf(w, "     to   : %s\n", result.DestPath)

	if len(result.Keys) == 0 {
		fmt.Fprintf(w, "     keys : (none)\n")
		return
	}

	sorted := make([]string, len(result.Keys))
	copy(sorted, result.Keys)
	sort.Strings(sorted)

	fmt.Fprintf(w, "     keys : %s\n", strings.Join(sorted, ", "))
}
