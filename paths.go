package graph

import (
	"fmt"
	"strings"
)

// Path is an ordered set of Nodes that make a path from the start,
// the first element in the slice, to the end, the last element in
// the slice.
type Path Nodes

// Paths is a collection of Path node sets.
type Paths []Path

// Identical checks if the given path is the same.
//
// Note: this currently uses the string representation, which might not always
//
//	be accurate if the nodes do not, or contain non-uniq names.
func (path Path) Identical(path2 Path) bool {
	return path.String() == path2.String()
}

// ContainsNode checks if the given node is part of the path.
func (path Path) ContainsNode(n *Node) bool {
	for _, pathNode := range path {
		if pathNode == n {
			return true
		}
	}
	return false
}

// String returns a human-readable string for the Path.
func (path Path) String() string {
	var builder strings.Builder

	for _, node := range path {
		builder.WriteString(fmt.Sprintf("→ %s ", node.Name))
	}

	return strings.TrimSpace(strings.TrimPrefix(builder.String(), "→ "))
}

// ContainsPath checks if the given path is identical to any of one
// of the path node sets.
func (paths Paths) ContainsPath(p Path) bool {
	for _, path := range paths {
		if path.Identical(p) {
			return true
		}
	}
	return false
}

// ContainsNode checks if the given node is contained in any one of
// the path node sets.
func (paths Paths) ContainsNode(n *Node) bool {
	for _, path := range paths {
		if path.ContainsNode(n) {
			return true
		}
	}
	return false
}
