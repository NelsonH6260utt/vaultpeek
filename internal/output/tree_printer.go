package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourusername/vaultpeek/internal/vault"
)

const (
	treeBranch = "├── "
	treeLast   = "└── "
	treeIndent = "│   "
	treeEmpty  = "    "
)

// PrintTree writes a visual tree of the given TreeNode to w.
func PrintTree(w io.Writer, node *vault.TreeNode) {
	fmt.Fprintln(w, node.Path)
	printChildren(w, node, "")
}

func printChildren(w io.Writer, node *vault.TreeNode, prefix string) {
	for i, child := range node.Children {
		isLast := i == len(node.Children)-1

		connector := treeBranch
		nextPrefix := prefix + treeIndent
		if isLast {
			connector = treeLast
			nextPrefix = prefix + treeEmpty
		}

		label := lastSegment(child.Path)
		fmt.Fprintf(w, "%s%s%s\n", prefix, connector, label)

		if !child.IsLeaf {
			printChildren(w, child, nextPrefix)
		}
	}
}

// lastSegment returns the final path component, preserving trailing slash
// for directories.
func lastSegment(path string) string {
	trailing := ""
	if strings.HasSuffix(path, "/") {
		trailing = "/"
		path = strings.TrimSuffix(path, "/")
	}
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return path
	}
	return parts[len(parts)-1] + trailing
}
