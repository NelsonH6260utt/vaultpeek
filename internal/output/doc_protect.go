// Package output provides formatting and printing utilities for vaultpeek
// command output.
//
// # Protect Printer
//
// PrintProtectResult renders the outcome of a [vault.ProtectSecret] call.
// A successful result prints the fully-qualified path together with a hint
// about the custom metadata key that was written. A failed result writes a
// diagnostic to stderr so that shell pipelines can distinguish between the
// two streams.
//
// Example output (success):
//
//	protected  secret/myapp/config
//	  The path is now marked with vaultpeek_protected=true in custom metadata.
//
// Example output (failure, stderr):
//
//	error: could not protect secret at "myapp/config"
//	  reason: permission denied
package output
