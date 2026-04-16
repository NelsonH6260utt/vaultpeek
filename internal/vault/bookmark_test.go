package vault

import (
	"testing"
)

func TestBookmarkResult_Fields(t *testing.T) {
	s := NewBookmarkStore()
	res := s.Add("prod-db", "secret", "prod/db/creds")
	if !res.Success {
		t.Fatal("expected success")
	}
	if res.Bookmark.Label != "prod-db" {
		t.Errorf("expected label prod-db, got %s", res.Bookmark.Label)
	}
	if res.Bookmark.Path != "prod/db/creds" {
		t.Errorf("unexpected path: %s", res.Bookmark.Path)
	}
	if !res.Added {
		t.Error("expected Added to be true")
	}
}

func TestBookmarkResult_FailureState(t *testing.T) {
	s := NewBookmarkStore()
	res := s.Add("", "secret", "prod/db")
	if res.Success {
		t.Error("expected failure for empty label")
	}
	if res.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestBookmarkStore_Remove(t *testing.T) {
	s := NewBookmarkStore()
	s.Add("mykey", "kv", "some/path")
	res := s.Remove("mykey")
	if !res.Success || !res.Removed {
		t.Error("expected successful removal")
	}
	_, ok := s.Get("mykey")
	if ok {
		t.Error("bookmark should have been deleted")
	}
}

func TestBookmarkStore_RemoveMissing(t *testing.T) {
	s := NewBookmarkStore()
	res := s.Remove("nonexistent")
	if res.Success {
		t.Error("expected failure removing missing bookmark")
	}
}

func TestBookmarkStore_ListSorted(t *testing.T) {
	s := NewBookmarkStore()
	s.Add("zebra", "kv", "z/path")
	s.Add("alpha", "kv", "a/path")
	s.Add("mango", "kv", "m/path")
	list := s.List()
	if len(list) != 3 {
		t.Fatalf("expected 3 bookmarks, got %d", len(list))
	}
	if list[0].Label != "alpha" || list[1].Label != "mango" || list[2].Label != "zebra" {
		t.Errorf("unexpected order: %v", list)
	}
}
