// Package output provides formatting and display utilities for vaultpeek.
//
// # Lock Printer
//
// PrintLockResult displays the outcome of a lock or unlock operation on a
// Vault secret path. Lock metadata is stored in KVv2 custom_metadata fields
// and can be inspected with the protect or audit commands.
//
// Output symbols:
//
//	[locked]   — secret has been marked as locked
//	[unlocked] — lock metadata has been cleared
//	[error]    — operation failed; error detail is shown inline
package output
