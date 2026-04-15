package vault

import (
	"testing"
)

func TestIsDir_WithSlash(t *testing.T) {
	if !IsDir("subpath/") {
		t.Error("expected IsDir to return true for key ending with '/'")
	}
}

func TestIsDir_WithoutSlash(t *testing.T) {
	if IsDir("mykey") {
		t.Error("expected IsDir to return false for key without trailing '/'")
	}
}

func TestIsDir_EmptyString(t *testing.T) {
	if IsDir("") {
		t.Error("expected IsDir to return false for empty string")
	}
}

func TestTrimDir_RemovesSlash(t *testing.T) {
	got := TrimDir("subpath/")
	if got != "subpath" {
		t.Errorf("TrimDir: expected %q, got %q", "subpath", got)
	}
}

func TestTrimDir_NoSlash(t *testing.T) {
	got := TrimDir("mykey")
	if got != "mykey" {
		t.Errorf("TrimDir: expected %q, got %q", "mykey", got)
	}
}

func TestTrimDir_EmptyString(t *testing.T) {
	got := TrimDir("")
	if got != "" {
		t.Errorf("TrimDir: expected empty string, got %q", got)
	}
}

func TestKVv2MetaPath_UsedByList(t *testing.T) {
	// Ensure the path helper produces the correct metadata prefix
	// which ListSecrets depends on.
	got := KVv2MetaPath("secret", "myapp/config")
	want := "secret/metadata/myapp/config"
	if got != want {
		t.Errorf("KVv2MetaPath: expected %q, got %q", want, got)
	}
}
