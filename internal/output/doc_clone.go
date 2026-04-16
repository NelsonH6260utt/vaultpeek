// Package output provides formatting and printing utilities for vaultpeek
// CLI output. This file documents the clone sub-package of output.
//
// Clone printing renders the result of a vault.CloneSecret operation,
// showing the source and destination paths along with a sorted list of
// keys that were copied. On failure the error message is surfaced with a
// clear failure indicator so the operator can diagnose the issue quickly.
package output
