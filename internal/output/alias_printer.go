package output

import (
	"fmt"
	"io"
	"os"

	"github.com/user/vaultpeek/internal/vault"
)

// PrintAliasResult writes the outcome of a set or remove alias operation.
func PrintAliasResult(r vault.AliasResult) {
	printAliasResultTo(os.Stdout, r)
}

func printAliasResultTo(w io.Writer, r vault.AliasResult) {
	if !r.Success {
		if r.Err != nil {
			fmt.Fprintf(w, "[error] alias operation failed: %v\n", r.Err)
		} else {
			fmt.Fprintln(w, "[error] alias operation failed")
		}
		return
	}
	if r.Path != "" {
		fmt.Fprintf(w, "[ok] alias %q -> %s\n", r.Name, r.Path)
	} else {
		fmt.Fprintf(w, "[ok] alias %q removed\n", r.Name)
	}
}

// PrintAliasList writes all stored aliases in tabular form.
func PrintAliasList(w io.Writer, entries []vault.AliasEntry) {
	if w == nil {
		w = os.Stdout
	}
	if len(entries) == 0 {
		fmt.Fprintln(w, "no aliases defined")
		return
	}
	maxLen := 0
	for _, e := range entries {
		if len(e.Name) > maxLen {
			maxLen = len(e.Name)
		}
	}
	fmt.Fprintln(w, "ALIAS\t\tPATH")
	fmt.Fprintln(w, repeatDash(40))
	for _, e := range entries {
		padding := maxLen - len(e.Name) + 2
		fmt.Fprintf(w, "%s%s%s\n", e.Name, spaces(padding), e.Path)
	}
}

func spaces(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = ' '
	}
	return string(s)
}
