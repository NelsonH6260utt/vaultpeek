// Package output provides formatting and printing utilities for vaultpeek CLI output.
//
// Tag Printer
//
// The tag printer renders the result of a TagSecret operation, showing
// the target path and all applied custom_metadata key-value pairs in
// sorted order. On failure, it emits an [error] prefixed message with
// the underlying error description.
package output
