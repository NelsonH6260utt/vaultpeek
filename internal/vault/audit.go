package vault

import (
	"fmt"
	"sort"
	"time"
)

// AuditEntry represents a single audit log record for a secret path.
type AuditEntry struct {
	Path      string
	Operation string
	Version   int
	Timestamp time.Time
	Renewable bool
}

// AuditLog holds a collection of audit entries for a given mount path.
type AuditLog struct {
	Mount   string
	Entries []AuditEntry
}

// SortByTime sorts audit entries in ascending chronological order.
func (a *AuditLog) SortByTime() {
	sort.Slice(a.Entries, func(i, j int) bool {
		return a.Entries[i].Timestamp.Before(a.Entries[j].Timestamp)
	})
}

// FilterByPath returns only entries matching the given secret path.
func (a *AuditLog) FilterByPath(path string) []AuditEntry {
	var result []AuditEntry
	for _, e := range a.Entries {
		if e.Path == path {
			result = append(result, e)
		}
	}
	return result
}

// FilterByOperation returns only entries matching the given operation type.
func (a *AuditLog) FilterByOperation(op string) []AuditEntry {
	var result []AuditEntry
	for _, e := range a.Entries {
		if e.Operation == op {
			result = append(result, e)
		}
	}
	return result
}

// Summary returns a human-readable summary string for the audit log.
func (a *AuditLog) Summary() string {
	return fmt.Sprintf("mount=%s entries=%d", a.Mount, len(a.Entries))
}
