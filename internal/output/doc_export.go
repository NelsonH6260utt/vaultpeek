// Package output provides formatting and display utilities for vaultpeek CLI output.
//
// The export_printer.go file in this package handles rendering of ExportResult
// values produced by vault.ExportSecret. It supports two modes:
//
//   - Human-readable table mode (default): displays key-value pairs sorted
//     alphabetically under the secret path header.
//
//   - JSON mode (--json flag): emits indented JSON suitable for piping to
//     other tools such as jq or redirecting to a file.
//
// Example usage:
//
//	result := vault.ExportSecret(client, "secret", "myapp/config")
//	output.PrintExportResult(result, jsonMode)
package output
