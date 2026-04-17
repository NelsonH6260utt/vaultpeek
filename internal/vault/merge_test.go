package vault

import (
	"errors"
	"testing"
)

func TestMergeResult_Fields(t *testing.T) {
	r := MergeResult{Path: "secret/foo", Keys: []string{"a", "b"}, Success: true}
	if r.Path != "secret/foo" {
		t.Errorf("expected path secret/foo, got %s", r.Path)
	}
	if len(r.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(r.Keys))
	}
}

func TestMergeResult_FailureState(t *testing.T) {
	r := MergeResult{Success: false, Err: errors.New("write failed")}
	if r.Success {
		t.Error("expected Success=false")
	}
	if r.Err == nil {
		t.Error("expected non-nil Err")
	}
}

func TestMergeSecret_KeepDest(t *testing.T) {
	src := map[string]interface{}{"a": "src-a", "b": "src-b"}
	dst := map[string]interface{}{"a": "dst-a"}
	var written map[string]interface{}
	write := func(_ string, data map[string]interface{}) error {
		written = data
		return nil
	}
	res := MergeSecret(src, dst, "secret/x", MergeStrategyKeepDest, write)
	if !res.Success {
		t.Fatal("expected success")
	}
	if written["a"] != "dst-a" {
		t.Errorf("expected dst-a kept, got %v", written["a"])
	}
	if written["b"] != "src-b" {
		t.Errorf("expected src-b added, got %v", written["b"])
	}
}

func TestMergeSecret_Overwrite(t *testing.T) {
	src := map[string]interface{}{"a": "src-a"}
	dst := map[string]interface{}{"a": "dst-a"}
	var written map[string]interface{}
	res := MergeSecret(src, dst, "secret/x", MergeStrategyOverwrite, func(_ string, data map[string]interface{}) error {
		written = data
		return nil
	})
	if !res.Success {
		t.Fatal("expected success")
	}
	if written["a"] != "src-a" {
		t.Errorf("expected src-a, got %v", written["a"])
	}
}

func TestMergeSecret_NilSource(t *testing.T) {
	res := MergeSecret(nil, map[string]interface{}{"a": "1"}, "secret/x", MergeStrategyKeepDest, func(_ string, _ map[string]interface{}) error { return nil })
	if res.Success {
		t.Error("expected failure on nil source")
	}
}

func TestMergeSecret_WriteError(t *testing.T) {
	res := MergeSecret(map[string]interface{}{"a": "1"}, nil, "secret/x", MergeStrategyKeepDest, func(_ string, _ map[string]interface{}) error {
		return errors.New("vault unavailable")
	})
	if res.Success {
		t.Error("expected failure on write error")
	}
	if res.Err == nil {
		t.Error("expected non-nil Err")
	}
}
