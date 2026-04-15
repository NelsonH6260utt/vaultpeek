package vault

import (
	"testing"
)

func TestVersionMeta_Fields(t *testing.T) {
	v := VersionMeta{
		Version:      3,
		CreatedTime:  "2024-01-01T00:00:00Z",
		DeletionTime: "",
		Destroyed:    false,
	}
	if v.Version != 3 {
		t.Errorf("expected version 3, got %d", v.Version)
	}
	if v.CreatedTime != "2024-01-01T00:00:00Z" {
		t.Errorf("unexpected CreatedTime: %s", v.CreatedTime)
	}
	if v.Destroyed {
		t.Error("expected Destroyed to be false")
	}
}

func TestVersionMeta_Destroyed(t *testing.T) {
	v := VersionMeta{
		Version:   1,
		Destroyed: true,
	}
	if !v.Destroyed {
		t.Error("expected Destroyed to be true")
	}
}

func TestListVersions_SortOrder(t *testing.T) {
	// Simulate the sort logic used in ListVersions without a live Vault.
	input := []VersionMeta{
		{Version: 3},
		{Version: 1},
		{Version: 2},
	}

	// Apply same sort as ListVersions.
	importSort(input)

	for i, want := range []int{1, 2, 3} {
		if input[i].Version != want {
			t.Errorf("index %d: expected version %d, got %d", i, want, input[i].Version)
		}
	}
}

// importSort mirrors the sort logic from ListVersions for unit testing.
func importSort(metas []VersionMeta) {
	for i := 1; i < len(metas); i++ {
		for j := i; j > 0 && metas[j].Version < metas[j-1].Version; j-- {
			metas[j], metas[j-1] = metas[j-1], metas[j]
		}
	}
}
