package graph

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
	Name      string
	Node      *Node
	Direction EdgeDirection
	Attributes
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

type AddEdge struct {
	From      *Node
	Direction *EdgeDirection
	To        *Node
}

func AddEdges(addEdge ...AddEdge) {
	for _, addEdge := range addEdge {
		if addEdge.Direction != nil {
			addEdge.From.AddEdgeWithDirection(addEdge.To, *addEdge.Direction)
			continue
		}
		addEdge.From.AddEdge(addEdge.To)
	}
}
