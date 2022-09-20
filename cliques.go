package graph


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