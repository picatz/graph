package graph


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
