package output

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/user/vaultpeek/internal/vault"
)

const (
	symDestroyed = "✖"
	symDeleted   = "⊘"
	symActive    = "✔"
)

// PrintVersionHistory writes a formatted table of secret version metadata to w.
func PrintVersionHistory(w io.Writer, path string, versions []vault.VersionMeta) {
	fmt.Fprintf(w, "Version history for: %s\n", path)
	fmt.Fprintln(w, repeatDash(60))

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "VERSION\tSTATUS\tCREATED\tDELETED")
	fmt.Fprintln(tw, "-------\t------\t-------\t-------")

	for _, v := range versions {
		status := symActive
		if v.Destroyed {
			status = symDestroyed
		} else if v.DeletionTime != "" {
			status = symDeleted
		}

		deletionDisplay := v.DeletionTime
		if deletionDisplay == "" {
			deletionDisplay = "-"
		}

		fmt.Fprintf(tw, "%d\t%s\t%s\t%s\n",
			v.Version,
			status,
			v.CreatedTime,
			deletionDisplay,
		)
	}
	tw.Flush()
}
