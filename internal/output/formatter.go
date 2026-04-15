// Package output provides formatting utilities for displaying
// Vault secret data and diff results in the terminal.
package output

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
)

// Format controls the output style.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes formatted output to a writer.
type Formatter struct {
	w      io.Writer
	format Format
}

// New creates a Formatter writing to w with the given format.
func New(w io.Writer, format Format) *Formatter {
	if w == nil {
		w = os.Stdout
	}
	return &Formatter{w: w, format: format}
}

// PrintSecrets renders a map of secret key/value pairs.
func (f *Formatter) PrintSecrets(path string, data map[string]interface{}) {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(f.w, "Path: %s\n", color.CyanString(path))
	fmt.Fprintf(f.w, "%s\n", strings.Repeat("-", 40))
	for _, k := range keys {
		fmt.Fprintf(f.w, "  %-24s = %v\n", color.YellowString(k), data[k])
	}
}

// PrintDiffLine writes a single diff line with appropriate coloring.
func (f *Formatter) PrintDiffLine(symbol, key string, value interface{}) {
	switch symbol {
	case "+":
		fmt.Fprintf(f.w, "%s %-24s = %v\n", color.GreenString("+"), key, value)
	case "-":
		fmt.Fprintf(f.w, "%s %-24s = %v\n", color.RedString("-"), key, value)
	case "~":
		fmt.Fprintf(f.w, "%s %-24s = %v\n", color.YellowString("~"), key, value)
	default:
		fmt.Fprintf(f.w, "  %-24s = %v\n", key, value)
	}
}

// PrintError writes an error message to the writer.
func (f *Formatter) PrintError(err error) {
	fmt.Fprintf(f.w, "%s %v\n", color.RedString("error:"), err)
}
