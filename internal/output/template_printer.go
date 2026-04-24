package output

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/user/vaultpeek/internal/vault"
)

// PrintTemplateResult writes a human-readable summary of a TemplateResult
// to stdout.
func PrintTemplateResult(r vault.TemplateResult) {
	printTemplateResultTo(os.Stdout, r)
}

func printTemplateResultTo(w io.Writer, r vault.TemplateResult) {
	if !r.Success {
		if r.Err != nil {
			fmt.Fprintf(w, "[error] template render failed for %q: %v\n", r.Path, r.Err)
			return
		}
		fmt.Fprintf(w, "[warn]  template rendered with %d unresolved placeholder(s) at %q\n",
			len(r.Missing), r.Path)
		for _, k := range r.Missing {
			fmt.Fprintf(w, "        missing: %s\n", k)
		}
		return
	}

	fmt.Fprintf(w, "[ok]    template rendered successfully for %q\n", r.Path)

	keys := make([]string, 0, len(r.Rendered))
	for k := range r.Rendered {
		if k == "__output" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "        %s = %s\n", k, r.Rendered[k])
	}

	if out, ok := r.Rendered["__output"]; ok {
		fmt.Fprintf(w, "--- rendered output ---\n%s\n", out)
	}
}
