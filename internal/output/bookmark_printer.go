package output

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/yourusername/vaultpeek/internal/vault"
)

// PrintBookmarkResult writes the result of an add/remove bookmark operation.
func PrintBookmarkResult(res vault.BookmarkResult) {
	printBookmarkResultTo(os.Stdout, res)
}

func printBookmarkResultTo(w io.Writer, res vault.BookmarkResult) {
	if !res.Success {
		fmt.Fprintf(w, "[error] bookmark operation failed: %v\n", res.Err)
		return
	}
	switch {
	case res.Added:
		fmt.Fprintf(w, "[+] bookmark added: %s => %s/%s\n", res.Bookmark.Label, res.Bookmark.Mount, res.Bookmark.Path)
	case res.Removed:
		fmt.Fprintf(w, "[-] bookmark removed: %s\n", res.Bookmark.Label)
	default:
		fmt.Fprintf(w, "[ok] bookmark operation completed\n")
	}
}

// PrintBookmarkList renders all bookmarks in a tabular format.
func PrintBookmarkList(bookmarks []vault.Bookmark) {
	printBookmarkListTo(os.Stdout, bookmarks)
}

func printBookmarkListTo(w io.Writer, bookmarks []vault.Bookmark) {
	if len(bookmarks) == 0 {
		fmt.Fprintln(w, "no bookmarks saved")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "LABEL\tMOUNT\tPATH")
	for _, b := range bookmarks {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", b.Label, b.Mount, b.Path)
	}
	tw.Flush()
}
