package graph

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
