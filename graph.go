package graph

import (
	"fmt"
	"strings"
)

// Node is the base unit of which graphs are formed.
type Node struct {
	Name   string
	Edges  Edges
	Weight float64
}

// EdgeDirection describes the "direction" of an edge relative
// to a node. A direction can be in one of three states, none,
// in, or out.
type EdgeDirection int

const (
	Unknown EdgeDirection = 0 // [ ┄ ] Edge has an unknown direction.
	None    EdgeDirection = 1 // [ - ] Edge has no direction, is undirected.
	In      EdgeDirection = 2 // [ ← ] Edge has inward direction.
	Out     EdgeDirection = 3 // [ → ] Edge has outward direction.
	Both    EdgeDirection = 4 // [ ↔ ] Edge has both inward and outward direction.
)

// Edge is a relationship with a Node, which can be directed if
// the edge is an "in" or "out" (directed) or neither (undirected).
type Edge struct {
	Node      *Node
	Direction EdgeDirection
	Magnitude float64
}

// Edges is a collection of Node relationships.
type Edges []*Edge

// AddEdge adds a directed relationship to a Node.
//
//   n → e
func (n *Node) AddEdge(e *Node) {
	n.Edges = append(n.Edges, &Edge{Node: e, Direction: Out})
	e.Edges = append(e.Edges, &Edge{Node: n, Direction: In})
}

// AddLink adds a bi-directional relationship to a Node, that
// is both inwards and outwards from either side.
//
//   n ↔ e
func (n *Node) AddLink(e *Node) {
	n.AddEdge(e)
	e.AddEdge(n)
}

// AddEdgeWithDirection adds a potentially directed relationship to a Node. The direction
// is up to the caller of the function. A corresponding edge is automatically added
// added; that is, if an "out edge" is added, an "in edge" is added on the other side
// of the edge relationship. This allows for the relationships to be bi-directionally
// walked from any point in the graph.
func (n *Node) AddEdgeWithDirection(e *Node, direction EdgeDirection) {
	switch direction {
	case None, Unknown, Both:
		n.Edges = append(n.Edges, &Edge{Node: e, Direction: direction})
		e.Edges = append(e.Edges, &Edge{Node: n, Direction: direction})
	case Out:
		n.Edges = append(n.Edges, &Edge{Node: e, Direction: Out})
		e.Edges = append(e.Edges, &Edge{Node: n, Direction: In})
	case In:
		n.Edges = append(n.Edges, &Edge{Node: e, Direction: In})
		e.Edges = append(e.Edges, &Edge{Node: n, Direction: Out})
	}
}

// HasCycles checks if the Node is part of a cycle. A cycle of a graph
// is a subset of the edge set of a graph that forms a path such that
// the first node of the path corresponds to the last.
//
// Example of Cycle
//
// a → b → c → a
//
// Example of Non-Cycle
//
// a → b → c
//
// https://mathworld.wolfram.com/GraphCycle.html
// https://en.wikipedia.org/wiki/Cycle_(graph_theory)
func (n *Node) HasCycles() bool {
	for _, edge := range n.Edges.Out() {
		if edge.Node.HasPath(n) {
			return true
		}
	}
	return false
}

// In returns the edges that are directed inwards (pointing to).
func (es Edges) In() Edges {
	var in Edges
	for _, e := range es {
		if e.Direction == In {
			in = append(in, e)
		}
	}
	return in
}

// In returns the edges that are directed outwards (pointing from).
func (es Edges) Out() Edges {
	var out Edges
	for _, e := range es {
		if e.Direction == Out {
			out = append(out, e)
		}
	}
	return out
}

// Visit walks the outward nodes with a depth-first algorithm.
func (n *Node) Visit(fn func(*Node)) {
	visit(n, nil, fn)
}

// VisitAll walks the the outwards and inwards nodes with a
// depth-first algorithm.
func (n *Node) VisitAll(fn func(*Node)) {
	visitAll(n, nil, fn)
}

// visitWithTerminator is an internal function used to walk node
// relationships starting at the root node.
func visitWithTerminator(root *Node, record map[*Node]struct{}, in, out bool, fn func(*Node) bool) {
	if root == nil {
		return
	}

	if record == nil {
		record = map[*Node]struct{}{}
	}

	_, alreadyVisited := record[root]
	if alreadyVisited {
		return
	}
	record[root] = struct{}{}

	if !fn(root) {
		return
	}

	for _, edge := range root.Edges {
		if out {
			if edge.Direction == Out {
				visitWithTerminator(edge.Node, record, in, out, fn)
			}
		}
		if in {
			if edge.Direction == In {
				visitWithTerminator(edge.Node, record, in, out, fn)
			}
		}
	}
}

// visit is an internal function that walks the outward nodes with
// a depth-first algorithm.
func visit(root *Node, record map[*Node]struct{}, fn func(*Node)) {
	wrapFn := func(n *Node) bool {
		fn(n)
		return true
	}

	visitWithTerminator(root, nil, false, true, wrapFn)
}

// visitAll is an internal function that walks the outward and inward
// nodes with a depth-first algorithm.
func visitAll(root *Node, record map[*Node]struct{}, fn func(*Node)) {
	wrapFn := func(n *Node) bool {
		fn(n)
		return true
	}

	visitWithTerminator(root, nil, true, true, wrapFn)
}

// Path is an ordered set of Nodes that make a path from the start,
// the first element in the slice, to the end, the last element in
// the slice.
type Path []*Node

// String returns a human-readable string for the Path.
func (path Path) String() string {
	var builder strings.Builder

	for _, node := range path {
		builder.WriteString(fmt.Sprintf("→ %s ", node.Name))
	}

	return strings.TrimSpace(strings.TrimPrefix(builder.String(), "→ "))
}

// PathTo returns the Path to the given end Node, nil if no path
// was found.
func (n *Node) PathTo(end *Node) Path {
	var hasPath bool
	var path Path

	visitWithTerminator(n, nil, false, true, func(n *Node) bool {
		path = append(path, n)
		for _, edge := range n.Edges.Out() {
			if edge.Node == end {
				path = append(path, edge.Node)
				hasPath = true
				return false // stop
			}
		}

		return true // continue
	})

	if !hasPath {
		return nil
	}

	return path
}

// HasPath checks if there is a Path to the given end Node.
//
//   root node       f           end node
//   ┌────────     ↗             ┌───────
//   a → b → c → e           i → e
//           ↓     ↘       ↗
//           d       g → h
//
//   Path: a → b → c → e → g → h → i → e
//
func (n *Node) HasPath(end *Node) bool {
	return n.PathTo(end) != nil
}

// ConnectNodes creats an ordered, directed relationship between
// the given nodes. The first node has an edge to the second node,
// which has a relationship to the third node, etc.
//
//   a → b → c → ...
//
func ConnectNodes(nodes ...*Node) {
	for i := range nodes {
		if i+1 < len(nodes) {
			x := nodes[i]
			y := nodes[i+1]
			x.AddEdge(y)
		}
	}
}

// MeshNodes creats a fully meshed, bi-directional relationship between
// all of the given nodes.
//
//       a
//    ⤢  ↑  ⤡
//   b ←─┼─→ d
//    ⤡  ↓  ⤢
//       c
//
func MeshNodes(nodes ...*Node) {
	for i := range nodes {
		if i+1 < len(nodes) {
			x := nodes[i]
			for _, y := range nodes[i+1:] {
				x.AddEdge(y)
				y.AddEdge(x)
			}
		}
	}
}

// FindBridges finds all "bridge" edges within a graph. An edge is a
// bridge if and only if it is not contained in any cycle. A bridge
// therefore cannot be a cycle chord.
//
// A "bridge" is also known as an "isthmus", "cut-edge", or "cut arc".
//
//          a ← d
//        ↙   ↖
//   e → b  →  c     Bridges (3): e → b, f → b, d → a
//       ↑
//       f
//
//   a           e
//   ↑ ⤡       ⤢ ↑
//   |   c → d   |   Bridges (1): c → d
//   ↓ ⤢       ⤡ ↓
//   b           f
//
//
// https://mathworld.wolfram.com/GraphBridge.html
// https://en.wikipedia.org/wiki/Bridge_(graph_theory)
func FindBridges(root *Node) []Path {
	bridges := []Path{}

	root.VisitAll(func(n *Node) {
		for _, edge := range n.Edges.Out() {
			if !n.HasCycles() {
				bridges = append(bridges, n.PathTo(edge.Node))
			} else {
				if !edge.Node.HasPath(n) {
					bridges = append(bridges, n.PathTo(edge.Node))
				}
			}
		}
	})

	return bridges
}
