package vault

import (
	"testing"
)

func TestTreeNode_Flatten_Leaf(t *testing.T) {
	node := &TreeNode{Path: "secret/foo", IsLeaf: true}
	paths := node.Flatten()
	if len(paths) != 1 {
		t.Fatalf("expected 1 path, got %d", len(paths))
	}
	if paths[0] != "secret/foo" {
		t.Errorf("expected 'secret/foo', got %q", paths[0])
	}
}

func TestTreeNode_Flatten_WithChildren(t *testing.T) {
	root := &TreeNode{
		Path: "secret/",
		Children: []*TreeNode{
			{Path: "secret/alpha", IsLeaf: true},
			{Path: "secret/beta", IsLeaf: true},
		},
	}
	paths := root.Flatten()
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	if paths[0] != "secret/alpha" {
		t.Errorf("expected 'secret/alpha', got %q", paths[0])
	}
	if paths[1] != "secret/beta" {
		t.Errorf("expected 'secret/beta', got %q", paths[1])
	}
}

func TestTreeNode_Flatten_Nested(t *testing.T) {
	root := &TreeNode{
		Path: "secret/",
		Children: []*TreeNode{
			{
				Path: "secret/dir/",
				Children: []*TreeNode{
					{Path: "secret/dir/key", IsLeaf: true},
				},
			},
			{Path: "secret/top", IsLeaf: true},
		},
	}
	paths := root.Flatten()
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d", len(paths))
	}
	if paths[0] != "secret/dir/key" {
		t.Errorf("expected 'secret/dir/key', got %q", paths[0])
	}
	if paths[1] != "secret/top" {
		t.Errorf("expected 'secret/top', got %q", paths[1])
	}
}

func TestTreeNode_Flatten_Empty(t *testing.T) {
	root := &TreeNode{Path: "secret/"}
	paths := root.Flatten()
	if len(paths) != 0 {
		t.Fatalf("expected 0 paths, got %d", len(paths))
	}
}
