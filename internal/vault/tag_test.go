package vault

import "testing"

func TestTagResult_Fields(t *testing.T) {
	r := TagResult{
		Path:    "secret/myapp",
		Tags:    map[string]string{"env": "prod", "owner": "team-a"},
		Success: true,
	}
	if r.Path != "secret/myapp" {
		t.Errorf("expected path secret/myapp, got %s", r.Path)
	}
	if len(r.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(r.Tags))
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestTagResult_FailureState(t *testing.T) {
	r := TagResult{Path: "secret/x", Success: false, Err: fmt.Errorf("boom")}
	if r.Success {
		t.Error("expected failure")
	}
	if r.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestSortedTagKeys_Order(t *testing.T) {
	tags := map[string]string{"zzz": "1", "aaa": "2", "mmm": "3"}
	keys := SortedTagKeys(tags)
	if keys[0] != "aaa" || keys[1] != "mmm" || keys[2] != "zzz" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestSortedTagKeys_Empty(t *testing.T) {
	keys := SortedTagKeys(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}

func TestTagResult_NilTags(t *testing.T) {
	r := TagResult{Path: "secret/empty"}
	if r.Tags != nil {
		t.Error("expected nil tags")
	}
	if r.Success {
		t.Error("expected Success false by default")
	}
}
