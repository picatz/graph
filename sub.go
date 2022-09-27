package graph

// Sub is a "subgraph", a subet of nodes within a graph. This can be
// used to namespace specific nodes so they're in logical graph.
//
// https://en.wikipedia.org/wiki/Induced_subgraph
type Sub = Instance

func NewSub(name string, attrs Attributes, nodes Nodes) *Sub {
	return New(name, attrs, nodes)
}
