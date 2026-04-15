package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/robbrockway/vaultpeek/internal/vault"
)

// PrintExportResult writes the export result to stdout in a human-readable
// or JSON format depending on the jsonMode flag.
func PrintExportResult(result vault.ExportResult, jsonMode bool) {
	printExportResultTo(os.Stdout, result, jsonMode)
}

func printExportResultTo(w io.Writer, result vault.ExportResult, jsonMode bool) {
	if !result.Success {
		fmt.Fprintf(w, "[error] export failed for %s: %v\n", result.Path, result.Err)
		return
	}

	if jsonMode {
		bytes, err := vault.MarshalExport(result)
		if err != nil {
			fmt.Fprintf(w, "[error] could not marshal export: %v\n", err)
			return
		}
		fmt.Fprintf(w, "%s\n", string(bytes))
		return
	}

	fmt.Fprintf(w, "Path: %s\n", result.Path)
	fmt.Fprintf(w, repeatDash(40)+"\n")

	keys := make([]string, 0, len(result.Data))
	for k := range result.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "  %-24s = %v\n", k, result.Data[k])
	}
}
