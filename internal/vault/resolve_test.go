package vault

import (
	"testing"
)

func TestResolveResult_Fields(t *testing.T) {
	r := ResolveResult{
		Path:     "myapp",
		Resolved: "secret/data/myapp",
		Success:  true,
	}
	if r.Path != "myapp" {
		t.Errorf("expected Path 'myapp', got %q", r.Path)
	}
	if r.Resolved != "secret/data/myapp" {
		t.Errorf("expected Resolved 'secret/data/myapp', got %q", r.Resolved)
	}
	if !r.Success {
		t.Error("expected Success to be true")
	}
}

func TestResolveSecret_EmptyPath(t *testing.T) {
	r := ResolveSecret("", ResolveOptions{})
	if r.Success {
		t.Error("expected failure for empty path")
	}
	if r.Err == nil {
		t.Error("expected non-nil Err for empty path")
	}
}

func TestResolveSecret_AliasMatch(t *testing.T) {
	opts := ResolveOptions{
		Aliases: map[string]string{
			"db": "secret/data/production/database",
		},
	}
	r := ResolveSecret("db", opts)
	if !r.Success {
		t.Fatalf("expected success, got err: %v", r.Err)
	}
	if r.Resolved != "secret/data/production/database" {
		t.Errorf("unexpected resolved path: %q", r.Resolved)
	}
}

func TestResolveSecret_NoAliasWithMount(t *testing.T) {
	opts := ResolveOptions{Mount: "secret/data"}
	r := ResolveSecret("myapp/config", opts)
	if !r.Success {
		t.Fatalf("expected success, got err: %v", r.Err)
	}
	if r.Resolved != "secret/data/myapp/config" {
		t.Errorf("unexpected resolved path: %q", r.Resolved)
	}
}

func TestResolveSecret_AlreadyMounted(t *testing.T) {
	opts := ResolveOptions{Mount: "secret/data"}
	r := ResolveSecret("secret/data/myapp", opts)
	if !r.Success {
		t.Fatalf("expected success")
	}
	if r.Resolved != "secret/data/myapp" {
		t.Errorf("path should not be double-prefixed, got %q", r.Resolved)
	}
}

func TestResolveSecret_NoMountNoAlias(t *testing.T) {
	r := ResolveSecret("plain/path", ResolveOptions{})
	if !r.Success {
		t.Fatalf("expected success")
	}
	if r.Resolved != "plain/path" {
		t.Errorf("expected path unchanged, got %q", r.Resolved)
	}
}
