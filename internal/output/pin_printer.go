package output

import (
	"fmt"
	"io"
	"os"

	"github.com/your-org/vaultpeek/internal/vault"
)

// PrintPinResult writes a human-readable summary of a PinResult to stdout.
func PrintPinResult(r vault.PinResult) {
	printPinResultTo(os.Stdout, r)
}

func printPinResultTo(w io.Writer, r vault.PinResult) {
	if !r.Success {
		errMsg := "unknown error"
		if r.Error != nil {
			errMsg = r.Error.Error()
		}
		fmt.Fprintf(w, "[pin] failed: %s\n", errMsg)
		return
	}

	fmt.Fprintf(w, "[pin] pinned %s @ version %d\n", r.Path, r.Version)
	fmt.Fprintf(w, "      pinned at: %s\n", r.PinnedAt.Format("2006-01-02 15:04:05 UTC"))
}
