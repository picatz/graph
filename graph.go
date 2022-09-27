package graph

// Instance describes a graph of zero or more nodes.
type Instance struct {
	Name string
	Attributes
	Nodes
}

// New returns a new instance of a graph.
func New(name string, attrs Attributes, nodes Nodes) *Instance {
	return &Instance{
		Name:       name,
		Nodes:      nodes,
		Attributes: attrs,
	}
}

// Visit walks the nodes of the graph.
//
// It does not perform depth-first-search, but the
// given function can choose to implement that using
// the node's Visit method.
func (inst *Instance) Visit(fn func(*Node)) {
	if fn == nil {
		return
	}

	for _, node := range inst.Nodes {
		fn(node)
	}
}

// IsAcyclic returns true if the nodes in the graph
// contains no cycles.
//
// https://mathworld.wolfram.com/AcyclicGraph.html
func (inst *Instance) IsAcyclic() bool {
	for _, node := range inst.Nodes {
		if node.HasCycles() {
			return false
		}
	}
	return true
}

// IsUnicyclic returns true if the graph contains
// only a single cycle.
//
// https://mathworld.wolfram.com/UnicyclicGraph.html
func (inst *Instance) IsUnicyclic() bool {
	var nCycles int
	for _, node := range inst.Nodes {
		if node.HasCycles() {
			nCycles++
			if nCycles > 1 {
				return false
			}
		}
	}
	return nCycles == 1
}

// IsBipartite returns true if the nodes in the graph
// is a Bipartite graph, also called a bigraph, where
// nodes can be decomposed into two disjoint sets such
// that no two nodes within the same set are adjacent.
//
// https://mathworld.wolfram.com/BipartiteGraph.html
func (inst *Instance) IsBipartite() bool {
	return inst.IsMultipartite(2)
}

// https://en.wikipedia.org/wiki/Multipartite_graph
func (inst *Instance) IsMultipartite(k int) bool {
	nodeSets := NodeSets{}

	for _, node := range inst.Nodes {
		// Handle the case where no node sets exist.
		if len(nodeSets) == 0 {
			nodeSets = append(nodeSets, NewNodeSet(node))
			continue
		}

		// Determine which node set the node should be
		// added to, based on its adjacency characteristics.
		targetSet, ok := nodeSets.GetSetNotAdjacentWith(node)
		if !ok {
			targetSet = NewNodeSet(node)
			nodeSets = append(nodeSets, targetSet)
			if len(nodeSets) > k {
				return false
			}
		} else {
			targetSet.Add(node)
		}
	}

	return len(nodeSets) == k
}
