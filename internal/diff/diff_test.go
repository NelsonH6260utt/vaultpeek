package diff_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpeek/internal/diff"
)

func TestCompare_Identical(t *testing.T) {
	left := map[string]interface{}{"key": "value", "foo": "bar"}
	right := map[string]interface{}{"key": "value", "foo": "bar"}

	result := diff.Compare(left, right)

	if result.HasDifferences() {
		t.Error("expected no differences, but found some")
	}
	if len(result.Identical) != 2 {
		t.Errorf("expected 2 identical keys, got %d", len(result.Identical))
	}
}

func TestCompare_OnlyInLeft(t *testing.T) {
	left := map[string]interface{}{"key": "value", "extra": "only-left"}
	right := map[string]interface{}{"key": "value"}

	result := diff.Compare(left, right)

	if !result.HasDifferences() {
		t.Error("expected differences")
	}
	if _, ok := result.OnlyInLeft["extra"]; !ok {
		t.Error("expected 'extra' to be in OnlyInLeft")
	}
}

func TestCompare_OnlyInRight(t *testing.T) {
	left := map[string]interface{}{"key": "value"}
	right := map[string]interface{}{"key": "value", "extra": "only-right"}

	result := diff.Compare(left, right)

	if _, ok := result.OnlyInRight["extra"]; !ok {
		t.Error("expected 'extra' to be in OnlyInRight")
	}
}

func TestCompare_Different(t *testing.T) {
	left := map[string]interface{}{"key": "old"}
	right := map[string]interface{}{"key": "new"}

	result := diff.Compare(left, right)

	pair, ok := result.Different["key"]
	if !ok {
		t.Fatal("expected 'key' to be in Different")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected diff values: %v", pair)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	left := map[string]interface{}{}
	right := map[string]interface{}{}

	result := diff.Compare(left, right)

	if result.HasDifferences() {
		t.Error("expected no differences for two empty maps")
	}
	if len(result.Identical) != 0 {
		t.Errorf("expected 0 identical keys, got %d", len(result.Identical))
	}
}

func TestSummary_NoDifferences(t *testing.T) {
	left := map[string]interface{}{"a": "1"}
	right := map[string]interface{}{"a": "1"}

	result := diff.Compare(left, right)
	summary := result.Summary("staging", "prod")

	if !strings.Contains(summary, "No differences found") {
		t.Errorf("expected no-diff message, got: %s", summary)
	}
}

func TestSummary_WithDifferences(t *testing.T) {
	left := map[string]interface{}{"key": "old", "gone": "yes"}
	right := map[string]interface{}{"key": "new", "added": "yes"}

	result := diff.Compare(left, right)
	summary := result.Summary("staging", "prod")

	if !strings.Contains(summary, "changed") {
		t.Errorf("expected changed entry in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "staging only") {
		t.Errorf("expected left-only entry in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "prod only") {
		t.Errorf("expected right-only entry in summary, got: %s", summary)
	}
}
