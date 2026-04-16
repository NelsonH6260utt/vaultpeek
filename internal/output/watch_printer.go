package output

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/user/vaultpeek/internal/vault"
)

// PrintWatchEvent writes a single watch poll result to stdout.
func PrintWatchEvent(result vault.WatchResult) {
	printWatchEventTo(os.Stdout, result)
}

func printWatchEventTo(w io.Writer, result vault.WatchResult) {
	timestamp := time.Now().Format("15:04:05")

	if result.Error != nil {
		fmt.Fprintf(w, "[%s] ERROR %s: %v\n", timestamp, result.Path, result.Error)
		return
	}

	if !result.Changed {
		fmt.Fprintf(w, "[%s] NO CHANGE %s (v%d)\n", timestamp, result.Path, result.Version)
		return
	}

	fmt.Fprintf(w, "[%s] CHANGED %s → version %d\n", timestamp, result.Path, result.Version)

	if len(result.Data) == 0 {
		fmt.Fprintf(w, "  (no data)\n")
		return
	}

	keys := make([]string, 0, len(result.Data))
	for k := range result.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "  %-24s = %v\n", k, result.Data[k])
	}
}
