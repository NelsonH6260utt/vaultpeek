package vault

import (
	"testing"
	"time"
)

func TestAnnotateResult_Fields(t *testing.T) {
	now := time.Now().UTC()
	r := AnnotateResult{
		Path:        "secret/myapp/config",
		Annotations: map[string]string{"owner": "team-a", "env": "prod"},
		UpdatedAt:   now,
		Success:     true,
	}

	if r.Path != "secret/myapp/config" {
		t.Errorf("expected path secret/myapp/config, got %s", r.Path)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
	if r.Err != nil {
		t.Errorf("expected no error, got %v", r.Err)
	}
	if len(r.Annotations) != 2 {
		t.Errorf("expected 2 annotations, got %d", len(r.Annotations))
	}
}

func TestAnnotateResult_FailureState(t *testing.T) {
	r := AnnotateResult{
		Path:    "secret/myapp/config",
		Success: false,
		Err:     errAnnotateFailed,
	}

	if r.Success {
		t.Error("expected Success to be false")
	}
	if r.Err == nil {
		t.Error("expected non-nil error")
	}
}

var errAnnotateFailed = annotateErr("annotation write failed")

type annotateErr string

func (e annotateErr) Error() string { return string(e) }

func TestSortedAnnotationKeys_Order(t *testing.T) {
	annotations := map[string]string{
		"zebra": "last",
		"alpha": "first",
		"mango": "middle",
	}
	keys := SortedAnnotationKeys(annotations)
	expected := []string{"alpha", "mango", "zebra"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %s, got %s", i, expected[i], k)
		}
	}
}

func TestSortedAnnotationKeys_Empty(t *testing.T) {
	keys := SortedAnnotationKeys(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %d keys", len(keys))
	}
}

func TestAnnotateResult_NilAnnotations(t *testing.T) {
	r := AnnotateResult{
		Path:        "secret/empty",
		Annotations: nil,
		Success:     true,
	}
	keys := SortedAnnotationKeys(r.Annotations)
	if len(keys) != 0 {
		t.Errorf("expected 0 keys for nil annotations, got %d", len(keys))
	}
}
