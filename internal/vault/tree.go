package vault

import (
	"context"
	"strings"
)

// TreeNode represents a node in a Vault secret path tree.
type TreeNode struct {
	Path     string
	Children []*TreeNode
	IsLeaf   bool
}

// BuildTree recursively lists all paths under the given root and builds
// a tree structure. mount is the KV v2 mount (e.g. "secret").
func BuildTree(ctx context.Context, c *Client, mount, root string) (*TreeNode, error) {
	node := &TreeNode{Path: root}

	entries, err := ListSecrets(ctx, c, mount, root)
	if err != nil {
		// If listing fails the path is a leaf secret, not a directory.
		node.IsLeaf = true
		return node, nil
	}

	for _, entry := range entries {
		childPath := strings.TrimSuffix(root, "/") + "/" + TrimDir(entry)
		if IsDir(entry) {
			child, err := BuildTree(ctx, c, mount, childPath+"/")
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, child)
		} else {
			node.Children = append(node.Children, &TreeNode{
				Path:   childPath,
				IsLeaf: true,
			})
		}
	}

	return node, nil
}

// Flatten returns all leaf paths in the tree in depth-first order.
func (n *TreeNode) Flatten() []string {
	if n.IsLeaf {
		return []string{n.Path}
	}
	var paths []string
	for _, child := range n.Children {
		paths = append(paths, child.Flatten()...)
	}
	return paths
}
