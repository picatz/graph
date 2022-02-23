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

// Edge is a relationship with a Node, which can be directed if
// the edge is an "in" or "out" (directed) or neither (undirected).
type Edge struct {
	*Node
	In        bool
	Out       bool
	Magnitude float64
}

// Edges is a collection of Node relationships.
type Edges []*Edge

// AddEdge adds a directed relationship to a Node.
//
// n → e
func (n *Node) AddEdge(e *Node) {
	n.Edges = append(n.Edges, &Edge{Node: e, Out: true})
	e.Edges = append(e.Edges, &Edge{Node: n, In: true})
}

// AddLink adds a bi-directional relationship to a Node, that
// is both inwards and outwards from either side.
//
// n ↔ e
func (n *Node) AddLink(e *Node) {
	n.AddEdge(e)
	e.AddEdge(n)
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
		if edge.HasPath(n) {
			return true
		}
	}
	return false
}

// In returns the edges that are directed inwards (pointing to).
func (es Edges) In() Edges {
	var in Edges
	for _, e := range es {
		if e.In {
			in = append(in, e)
		}
	}
	return in
}

// In returns the edges that are directed outwards (pointing from).
func (es Edges) Out() Edges {
	var out Edges
	for _, e := range es {
		if e.Out {
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
			if edge.Out {
				visitWithTerminator(edge.Node, record, in, out, fn)
			}
		}
		if in {
			if edge.In {
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

	// if len(path) == 1 && n == end {
	// 	return nil
	// }

	return path
}

// HasPath checks if there is a Path to the given end Node.
func (n *Node) HasPath(end *Node) bool {
	return n.PathTo(end) != nil
}

func (n *Node) PathToWithout(end, without *Node) Path {
	// TODO
	return nil
}

func (n *Node) HasPathToWithout(end, without *Node) bool {
	return n.PathToWithout(end, without) != nil
}

// ConnectNodes creats an ordered, directed relationship between
// the given nodes. The first node has an edge to the second node,
// which has a relationship to the third node, etc.
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

// FindBridges finds all "bridge" edges within a connected graph. An edge is a
// bridge if and only if it is not contained in any cycle. A bridge therefore
// cannot be a cycle chord.
//
// A "bridge" is also known as an "isthmus", "cut-edge", or "cut arc".
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
				path := edge.PathTo(n)
				if len(path) == 0 {
					bridges = append(bridges, n.PathTo(edge.Node))
				} else if len(path) == 2 {
					// TODO, maybe check if the path[1] can go back to
					// edges from path[0] maybe to prove
					// its deletion would cause a break?

					// Maybe "HasPathToWithout" variant?
				}
			}
		}
	})

	return bridges
}
