package graph

import (
	"fmt"
	"sort"
	"strings"
)

// Sub is a "subgraph", a subet of nodes within a graph. This can be
// used to namespace specific nodes so they're in logical graph.
//
// https://en.wikipedia.org/wiki/Induced_subgraph
type Sub struct {
	Name  string
	Nodes NodeSet
	Root  *Node
	Attributes
}

// Visit walks the nodes of the subraph.
//
// It does not perform  a depth-first-search, but the
// given function can choose to implement that using
// the node's Visit method.
func (s *Sub) Visit(fn func(*Node)) {
	if fn == nil {
		return
	}

	for node := range s.Nodes {
		fn(node)
	}
}

// Attributes are named values that can be associated with a node or subgraph.
type Attributes map[string]any

// UseAttribute is a helper function to use a named attribute of a specific type.
func UseAttribute[T any](attrs Attributes, name string, fn func(T)) error {
	v, ok := attrs[name]
	if !ok {
		return fmt.Errorf("graph node attribute %q doesn't exist", name)
	}
	vt, ok := v.(T)
	if !ok {
		return fmt.Errorf("graph node attribute %q is of type %T not %T", name, v, vt)
	}
	fn(vt)
	return nil
}

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

// EdgeDirection describes the "direction" of an edge relative
// to a node. A direction can be in one of five states:
//  0. Unknown
//  1. None
//  2. In
//  3. Out
//  4. Both
type EdgeDirection int

const (
	Unknown EdgeDirection = 0 // [ ┄ ] Edge has an unknown direction.
	None    EdgeDirection = 1 // [ - ] Edge has no direction, is undirected.
	In      EdgeDirection = 2 // [ ← ] Edge has inward direction.
	Out     EdgeDirection = 3 // [ → ] Edge has outward direction.
	Both    EdgeDirection = 4 // [ ↔ ] Edge has both inward and outward direction.
)

// String returns a human and command-line friendly representation of the edge direction.
func (d EdgeDirection) String() string {
	switch d {
	case None:
		return "-"
	case In:
		return "←"
	case Out:
		return "→"
	case Both:
		return "↔"
	default: // Unknown and anything else.
		return "┄"
	}
}

// AnyOf checks if the edge direction is any of the given directions.
func (d EdgeDirection) AnyOf(directions ...EdgeDirection) bool {
	for _, direction := range directions {
		if direction == d {
			return true
		}
	}

	return false
}

// Match checks if the edge direction matches the given edge direction,
// either exactly, or implicitly.
//
//  0. Unknown ┄ : ┄
//  1. None    - : -, ↔
//  2. In      ← : ←, ↔
//  3. Out     → : →, ↔
//  4. Both    ↔ : ↔, ←, →
func (d EdgeDirection) Match(direction EdgeDirection) bool {
	switch direction {
	case None:
		return d == None
	case In:
		return d.AnyOf(In, Both)
	case Out:
		return d.AnyOf(Out, Both)
	case Both:
		return d.AnyOf(In, Out, Both)
	default:
		return d == Unknown || d > 4
	}
}

// Edge is a relationship with a Node, which can be directed if
// the edge is an "in" or "out" (directed) or neither (undirected).
type Edge struct {
	Node      *Node
	Direction EdgeDirection
	Magnitude float64
}

// Edges is a collection of Node relationships.
type Edges []*Edge

func (edges Edges) Contains(n *Node) bool {
	for _, edge := range edges {
		if edge.Node == n {
			return true
		}
	}
	return false
}

func (edges Edges) ButNotWith(n *Node) Edges {
	other := Edges{}
	for _, edge := range edges {
		if edge.Node != n {
			other = append(other, edge)
		}
	}
	return other
}

func (edges Edges) AdjacentNodes() NodeSet {
	nodeSet := NodeSet{}

	for _, edge := range edges {
		nodeSet.Add(edge.Node)
	}

	return nodeSet
}

func (edges Edges) AdjacentTo(nodes ...*Node) bool {
	nodeSet := NodeSet{}

	for _, edge := range edges {
		for _, node := range nodes {
			if edge.Node == node {
				nodeSet.Add(node)
			}
		}
	}

	return len(nodeSet) == len(nodes)
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

// FindBridges finds all "bridge" paths within a graph. An edge,
// part of a path, is a bridge if and only if it is not contained
// in any cycle. Therefore, a bridge cannot be a cycle chord.
//
// A "bridge" is also known as an "isthmus", "cut-edge", or "cut arc".
//
//	       a ← d
//	     ↙   ↖
//	e → b  →  c     Bridges (3): e → b, f → b, d → a
//	    ↑
//	    f
//
//	a           e
//	↑ ⤡       ⤢ ↑
//	|   c → d   |   Bridges (1): c → d
//	↓ ⤢       ⤡ ↓
//	b           f
//
// To find the bridges in a graph, we need to visit each node
// and determine if it contains an edge that, if removed, would
// disconnect the graph into two. This is, if the number of
// components increases.
//
// A bridge, isthmus, cut-edge, or cut arc is an edge of a
// graph whose deletion increases the graph's number of
// connected components. Equivalently, an edge is a bridge
// if and only if it is not contained in any cycle.
//
// References
// - https://en.wikipedia.org/wiki/Bridge_(graph_theory)
// - https://en.wikipedia.org/wiki/Strongly_connected_component
// - https://mathworld.wolfram.com/GraphBridge.html
func FindBridges(root *Node) []Path {
	bridges := Paths{}

	var addUniqBridge = func(p Path) {
		if len(p) == 0 {
			return
		}
		if !bridges.ContainsPath(p) {
			bridges = append(bridges, p)
		}
	}

	root.VisitAll(func(n *Node) {
		for _, edge := range n.Edges {
			// First, skip edge nodes that themselves do not contain edges.
			if len(edge.Node.Edges) == 0 {
				continue
			}

			// Second, handle the simple case of dangling edges. This is only
			// useful for simple cases, but avoids using more complex traversal
			// until it is actually needed, making the algorithm a bit simpler
			// to digest, because you can do so in distinct steps.
			//
			// Graph
			//
			//        a ← d
			//      ↙   ↖
			// e → b  →  c
			//     ↑
			//     f
			//
			// Bridges
			//
			// 1. e → b
			// 2. f → b
			// 3. d → a
			//
			// Cycles
			//
			// 1. a → b → c → a
			//

			if len(edge.Node.Edges) == 1 {
				path := edge.Node.PathTo(edge.Node.Edges[0].Node)
				if len(path) > 0 {
					addUniqBridge(path)
				}
				continue
			}

			// Third, we must be dealing with a non-simple case.
			//
			// Graph
			//
			//   edgeNodeEdge.Node.Edge[0]
			//            |
			// n          d
			// |        ↗   ↘
			// a → b → c  ←  e
			//     |   |
			// edge.Node
			//         |
			// edgeNodeEdge.Node
			//
			// Bridgs
			//
			// 1. a → b
			// 2. b → c
			//
			// Cycles
			//
			// 1. c → d → e → c
			//

			for _, edgeNodeEdge := range edge.Node.Edges {
				if !edgeNodeEdge.Node.HasPath(edge.Node) {
					path := edge.Node.PathTo(edgeNodeEdge.Node)
					if len(path) > 0 {
						addUniqBridge(path)
						continue
					}
				}

				// The edge direction might be Both which is not
				// currently handled by this function...
				//
				// Started hacking around with what they might look like,
				// but have no tests to confirm it works, or not:
				//
				// if edgeNodeEdge.Node == n {
				// 	continue
				// }
				//
				// if edgeNodeEdge.Direction == Both {
				// 	if len(edgeNodeEdge.Node.Edges) == 1 {
				// 		path := edge.Node.PathTo(edgeNodeEdge.Node)
				// 		if len(path) > 0 {
				// 			addUniqBridge(path)
				// 			continue
				// 		}
				// 	}
				//
				// 	for _, edgeNodeEdgeNodeEdge := range edgeNodeEdge.Node.Edges {
				// 		if edgeNodeEdgeNodeEdge.Node == edge.Node {
				// 			continue // skip
				// 		}
				// 		if !edge.Node.PathToWithout(edgeNodeEdgeNodeEdge.Node, edgeNodeEdge.Node) {
				// 			path := edge.Node.PathTo(edgeNodeEdge.Node)
				// 			if len(path) > 0 {
				// 				addUniqBridge(path)
				// 				continue
				// 			}
				// 		}
				// 	}
				// }
			}

			// Another useful example to consider while you're here:
			//
			//
			//     edge.Node
			//         |
			//     n   |    edgeNodeEdge.Node.Edge[0]
			//     |   |       |
			//     a   |       e
			//     ↑ ⤡ |     ⤢ ↑
			//     |   c ↔ d   |
			//     ↓ ⤢     | ⤡ ↓
			//     b       |   f
			//             |
			//             |
			// edgeNodeEdge.Node
			//
		}
	})

	return bridges
}

// Clique is a subset of nodes in a graph such that every two
// distinct nodes in the set are adjacent.
//
// https://en.wikipedia.org/wiki/Clique_(graph_theory)
type Clique = NodeSet

// Cliques is a collection of clique node sets.
type Cliques []Clique

func (cliques Cliques) ContainsClique(c Clique) bool {
	for _, clique := range cliques {
		if clique.SameAs(c) {
			return true
		}
	}
	return false
}

func (cliques Cliques) ContainsNode(n *Node) bool {
	for _, clique := range cliques {
		if clique.Contains(n) {
			return true
		}
	}
	return false
}

func (cliques Cliques) ContainsNodeWithIndex(n *Node) (int, bool) {
	for index, clique := range cliques {
		if clique.Contains(n) {
			return index, true
		}
	}
	return 0, false
}

// FindCliques handles finding all "cliques" within a graph. A a clique
// is a subset of nodes in a graph such that every two distinct nodes
// in the clique are adjacent. That is, a clique of a graph "G" is an
// induced subgraph of "G" that is complete.
//
// References
// - https://en.wikipedia.org/wiki/Clique_(graph_theory)
// - https://en.wikipedia.org/wiki/Induced_subgraph
// - https://en.wikipedia.org/wiki/Complete_graph
// - https://mathworld.wolfram.com/Clique.html
func FindCliques(root *Node, minSize int) Cliques {
	cliques := Cliques{}

	//           b
	//         ↙   ↖
	//       c       a
	//     ↙   ↘   ↗
	//    e  →   d
	//
	//
	// Cliques: [1] {c, e, d}

	root.VisitAll(func(n *Node) {
		if len(n.Edges) == 0 {
			return
		}

		clique := Clique{}
		clique.Add(n)

		for _, edge := range n.Edges {
			for _, otherEdge := range n.Edges.ButNotWith(edge.Node) {
				if otherEdge.Node.Edges.AdjacentTo(clique.Nodes()...) {
					clique.Add(otherEdge.Node)
				}
			}
		}

		if len(clique) >= minSize && !cliques.ContainsClique(clique) {
			cliques = append(cliques, clique)
		}
	})

	// Basically a tree structure...
	// groups := map[*Node]NodeSet{}
	// visitAll(root, nil, func(n *Node) {
	// 	fmt.Println(n.Name)
	// 	_, ok := groups[n]
	// 	if !ok {
	// 		groups[n] = NodeSet{}
	// 	}
	// 	for _, edge := range n.Edges {
	// 		groups[n][edge.Node] = struct{}{}
	// 	}
	// })

	return cliques
}
