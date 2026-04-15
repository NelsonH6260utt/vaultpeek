package output

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/vaultpeek/internal/vault"
)

// PrintSearchResults writes search results to w in a human-readable table.
func PrintSearchResults(w io.Writer, results []vault.SearchResult, query string) {
	if len(results) == 0 {
		fmt.Fprintf(w, "No results found for %q\n", query)
		return
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Path != results[j].Path {
			return results[i].Path < results[j].Path
		}
		return results[i].MatchedKey < results[j].MatchedKey
	})

	fmt.Fprintf(w, "Search results for %q (%d match(es)):\n\n", query, len(results))

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "PATH\tMATCHED KEY")
	fmt.Fprintln(tw, repeatDash(40)+"\t"+repeatDash(20))
	for _, r := range results {
		fmt.Fprintf(tw, "%s\t%s\n", r.Path, r.MatchedKey)
	}
	_ = tw.Flush()
}
