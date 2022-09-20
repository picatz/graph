package graph

import (
	"fmt"
	"sort"
	"strings"
)

// Node is the base unit of which graphs are formed.
type Node struct {
	// The (unique) name (or label).
	Name string
	// Adjacency list of edges.
	Edges Edges
	// Named attributes about the node.
	Attributes
}

// NewNode returns a new node with the given name and attributes.
func NewNode(name string, attrs Attributes) *Node {
	return &Node{
		Name:       name,
		Attributes: attrs,
	}
}

// Nodes is a collection of Node objects.
type Nodes []*Node

func (nodes Nodes) Names() []string {
	names := make([]string, len(nodes))

	for i, node := range nodes {
		names[i] = node.Name
	}

	return names
}

func (nodes Nodes) String() string {
	return strings.Join(nodes.Names(), ", ")
}

// NodeSet is a collection of uniqe Node objects. Meant to be useful for
// algorithms that require collections of nodes that should not have
// repeated sequences.
//
// Also particularly useful for recording visited nodes
// during graph traversal.
type NodeSet map[*Node]struct{}

// NewNodeSet returns a new NodeSet that includes the given nodes.
func NewNodeSet(nodes ...*Node) NodeSet {
	ns := NodeSet{}

	for _, n := range nodes {
		ns[n] = struct{}{}
	}

	return ns
}

func (n NodeSet) String() string {
	nodes := []string{}

	for node := range n {
		nodes = append(nodes, node.Name)
	}

	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i] < nodes[j]
	})

	return strings.Join(nodes, ", ")
}

func (n NodeSet) Contains(node *Node) bool {
	if len(n) == 0 {
		return false
	}
	_, ok := n[node]
	return ok
}

func (ns NodeSet) Add(node *Node) {
	ns[node] = struct{}{}
}

func (ns NodeSet) Nodes() []*Node {
	nodes := []*Node{}
	for n := range ns {
		nodes = append(nodes, n)
	}
	return nodes
}

func (ns NodeSet) SameAs(other NodeSet) bool {
	if len(ns) != len(other) {
		return false
	}

	var sameCount int

	for n := range ns {
		if _, ok := other[n]; !ok {
			return false
		}
		sameCount++
	}

	return len(ns) == sameCount
}

func (nodes Nodes) IndexOf(o *Node) int {
	for i, node := range nodes {
		if node == o {
			return i
		}
	}
	return -1
}

func (nodes Nodes) AtIndex(i int) (*Node, error) {
	if len(nodes) <= i {
		return nil, fmt.Errorf("graph invalid index %d for nodes of size %d", i, len(nodes))
	}
	return nodes[i], nil
}

// AddEdge adds a directed relationship to a Node.
//
//	n → e
//
// To control the direction used for the relationship, use the AddEdgeWithDirection method.
func (n *Node) AddEdge(e *Node) {
	n.Edges = append(n.Edges, &Edge{Node: e, Direction: Out})
	e.Edges = append(e.Edges, &Edge{Node: n, Direction: In})
}

// AddLink adds a bi-directional relationship to a Node.
//
// Note: while this is sometimes rendered with a single "↔" (Both),
//
//	    this method really defines two distinct edges using the
//	    In and Out direction.
//
//	n ↔ e : [ n → e, e → n ]
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
// # Example of Cycle
//
// a → b → c → a
//
// # Example of Non-Cycle
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

func (n *Node) HasCycleContaining(node *Node) bool {
	for _, edge := range n.Edges.Out() {
		path := edge.Node.PathTo(n)
		if len(path) > 0 {
			if path.ContainsNode(node) {
				return true
			}
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

// Visit walks the outward nodes using a depth-first-search.
//
//	   root node
//	   ┌────────         1. Start at root "a"
//	1  a           e 5   2. Go to edge node "b"
//	   ↑ ⤡ 3   4 ⤢ ↑     3. Go to edge node "c"
//	   |   c ↔ d   |     4. Go to edge node "d"
//	   ↓ ⤢       ⤡ ↓     5. Go to edge node "e"
//	2  b           f 6   6. Go to edge node "f"
func (n *Node) Visit(fn func(*Node)) {
	visit(n, nil, fn)
}

// VisitAll walks the the outwards and inwards nodes with a
// depth-first-search algorithm.
func (n *Node) VisitAll(fn func(*Node)) {
	visitAll(n, nil, fn)
}

// visitWithTerminator is an internal function used to walk node
// relationships starting at the root node using depth-first-search.
//
// The record node set keeps track of nodes which were already visited,
// to prevent infinite loops that can be found during traversal. The first
// call to this function can provide a nil record.
//
// The direction defines the edges which should be visted: "out" to walk
// outward edges, "in" to walk inward edge; "unknown", "none",
// and "both" can all be used to walk bi-directionally.
//
// Lastly, the function given to run for each visited node can return true
// to continue traversal, or false to stop traversal.
func visitWithTerminator(root *Node, record NodeSet, direction EdgeDirection, fn func(*Node) bool) {
	if root == nil {
		return
	}

	if record == nil {
		record = NodeSet{}
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
		switch direction {
		case Unknown, None, Both:
			visitWithTerminator(edge.Node, record, direction, fn)
		case In, Out:
			if edge.Direction == direction || edge.Direction == Both {
				visitWithTerminator(edge.Node, record, direction, fn)
			}
		}
	}
}

// visit is an internal function that walks the outward nodes with
// a depth-first algorithm.
func visit(root *Node, record NodeSet, fn func(*Node)) {
	wrapFn := func(n *Node) bool {
		fn(n)
		return true
	}

	visitWithTerminator(root, nil, Out, wrapFn)
}

// visitAll is an internal function that walks the outward and inward
// nodes with a depth-first algorithm.
func visitAll(root *Node, record NodeSet, fn func(*Node)) {
	wrapFn := func(n *Node) bool {
		fn(n)
		return true
	}

	visitWithTerminator(root, nil, Both, wrapFn)
}

// PathTo returns the Path to the given end Node, nil if no path
// was found.
func (n *Node) PathTo(end *Node) Path {
	var (
		hasPath bool
		path    Path
	)

	visitWithTerminator(n, nil, Out, func(n *Node) bool {
		if hasPath {
			return false
		}

		path = append(path, n)

		for _, edge := range n.Edges {
			switch edge.Direction {
			case Out, Both, None, Unknown:
				if edge.Node == end {
					path = append(path, edge.Node)
					hasPath = true
					return false
				}
			}
		}

		return !hasPath
	})

	if !hasPath {
		return nil
	}

	return path
}

// PathToWithout checks if there's a path to the given end node, without
// having to "go through" or "use" the other given node.
func (n *Node) PathToWithout(end, without *Node) bool {
	path := n.PathTo(end)
	return !path.ContainsNode(without)
}

// HasPath checks if there is a Path to the given end Node.
//
//	root node       f           end node
//	┌────────     ↗             ┌───────
//	a → b → c → e           i → e
//	        ↓     ↘       ↗
//	        d       g → h
//
//	Path: a → b → c → e → g → h → i → e
func (n *Node) HasPath(end *Node) bool {
	return n.PathTo(end) != nil
}

// ConnectNodes creats an ordered, directed relationship between
// the given nodes. The first node has an edge to the second node,
// which has a relationship to the third node, etc.
//
//	a → b → c → ...
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
//	    a
//	 ⤢  ↑  ⤡
//	b ←─┼─→ d
//	 ⤡  ↓  ⤢
//	    c
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
