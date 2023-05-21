package graph

// Instance describes a graph of zero or more nodes.
type Instance struct {
	// Name is the name of the graph instance.
	Name string

	// Attributes is a map of key-value pairs that describe the graph instance.
	Attributes

	// Nodes is a slice of nodes that belong to the graph instance.
	Nodes
}

// WithAttributes is a functional option that sets the attributes of the graph.
func WithAttributes(attrs Attributes) func(*Instance) {
	return func(inst *Instance) {
		inst.Attributes = attrs
	}
}

// WithNodes is a functional option that sets the nodes of the graph.
func WithNodes(nodes Nodes) func(*Instance) {
	return func(inst *Instance) {
		inst.Nodes = nodes
	}
}

// New returns a new instance of a graph.
func New(name string, opts ...func(*Instance)) *Instance {
	inst := &Instance{
		Name:       name,
		Nodes:      Nodes{},
		Attributes: Attributes{},
	}

	for _, opt := range opts {
		opt(inst)
	}

	return inst
}

// AddNode adds a node to the graph.
func (inst *Instance) AddNode(node *Node) {
	if node == nil {
		return
	}

	inst.Nodes = append(inst.Nodes, node)
}

// AddNodes adds a slice of nodes to the graph.
func (inst *Instance) AddNodes(nodes ...*Node) {
	if nodes == nil {
		return
	}

	inst.Nodes = append(inst.Nodes, nodes...)
}

// AddEdge adds an edge to the graph from the source node to the target node.
func (inst *Instance) AddEdge(from, to *Node) {
	if from == nil || to == nil {
		return
	}

	from.AddEdge(to)
}

// AddEdges adds a slice of edges to the graph.
func (inst *Instance) AddEdges(em EdgeMap) {
	for from, to := range em {
		for _, to := range to {
			inst.AddEdge(from, to)
		}
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

// DFS performs a depth-first-search of the graph.
//
// https://en.wikipedia.org/wiki/Depth-first_search
func (inst *Instance) DFS(fn func(*Node)) {
	if fn == nil {
		return
	}

	// Create a map of nodes that have been visited.
	visited := NodeSet{}

	// Iterate over all the nodes in the graph.
	for _, node := range inst.Nodes {
		// If the node has already been visited, skip it.
		if visited.Contains(node) {
			continue
		}

		// Create a stack of nodes to visit.
		stack := Nodes{}

		// Add the node to the stack.
		stack = append(stack, node)

		// While there are nodes in the stack, visit them.
		for len(stack) > 0 {
			// Get the last node in the stack.
			node := stack[len(stack)-1]

			// Remove the node from the stack.
			stack = stack[:len(stack)-1]

			// If the node has already been visited, skip it.
			if visited.Contains(node) {
				continue
			}

			// Visit the node.
			fn(node)

			// Mark the node as visited.
			visited.Add(node)

			// Add the node's children to the stack.
			stack = append(stack, node.Out().Nodes()...)
		}
	}
}

// BFS performs a breadth-first-search of the graph.
//
// https://en.wikipedia.org/wiki/Breadth-first_search
func (inst *Instance) BFS(fn func(*Node)) {
	if fn == nil {
		return
	}

	// Create a map of nodes that have been visited.
	visited := NodeSet{}

	// Iterate over all the nodes in the graph.
	for _, node := range inst.Nodes {
		// If the node has already been visited, skip it.
		if visited.Contains(node) {
			continue
		}

		// Create a queue of nodes to visit.
		queue := Nodes{}

		// Add the node to the queue.
		queue = append(queue, node)

		// While there are nodes in the queue, visit them.
		for len(queue) > 0 {
			// Get the first node in the queue.
			node := queue[0]

			// Remove the node from the queue.
			queue = queue[1:]

			// If the node has already been visited, skip it.
			if visited.Contains(node) {
				continue
			}

			// Visit the node.
			fn(node)

			// Mark the node as visited.
			visited.Add(node)

			// Add the node's children to the queue.
			queue = append(queue, node.Out().Nodes()...)
		}
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
