package vault

import (
	"testing"
	"time"
)

func TestArchiveResult_Fields(t *testing.T) {
	r := ArchiveResult{
		Path:    "secret/myapp/config",
		Success: true,
		Keys:    []string{"host", "port"},
	}
	if r.Path != "secret/myapp/config" {
		t.Errorf("unexpected path: %s", r.Path)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if len(r.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Keys))
	}
}

func TestArchiveResult_FailureState(t *testing.T) {
	r := ArchiveResult{
		Path:    "secret/missing",
		Success: false,
	}
	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Error != nil {
		t.Errorf("expected nil error, got %v", r.Error)
	}
}

func TestArchiveEntry_Fields(t *testing.T) {
	now := time.Now().UTC()
	e := ArchiveEntry{
		Path:       "secret/app/db",
		Data:       map[string]interface{}{"password": "s3cr3t"},
		ArchivedAt: now,
	}
	if e.Path != "secret/app/db" {
		t.Errorf("unexpected path: %s", e.Path)
	}
	if e.Data["password"] != "s3cr3t" {
		t.Errorf("unexpected data value")
	}
	if e.ArchivedAt.IsZero() {
		t.Error("ArchivedAt should not be zero")
	}
}

func TestArchiveEntry_EmptyData(t *testing.T) {
	e := ArchiveEntry{
		Path: "secret/empty",
		Data: map[string]interface{}{},
	}
	if len(e.Data) != 0 {
		t.Errorf("expected empty data map")
	}
}

func TestArchiveResult_NilKeys(t *testing.T) {
	r := ArchiveResult{}
	if r.Keys != nil {
		t.Errorf("expected nil keys slice")
	}
}
