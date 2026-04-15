// Package output provides terminal formatting for vaultpeek.
//
// It supports rendering secret key/value pairs and diff results
// with ANSI color coding to improve readability in the terminal.
//
// Supported output formats:
//
//	FormatText - human-readable colored terminal output
//	FormatJSON - structured JSON output (future)
//
// Example usage:
//
//	f := output.New(os.Stdout, output.FormatText)
//	f.PrintSecrets("secret/myapp/prod", data)
package output
