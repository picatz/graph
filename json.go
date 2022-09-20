package graph

import (
	"encoding/json"
	"fmt"
	"io"
)

type nodeJSON struct {
	Name       string `json:"name,omitempty"`
	Attributes `json:"attributes,omitempty"`
}

type edgeJSON struct {
	Name       string        `json:"name,omitempty"`
	FromIndex  int           `json:"from_index"`
	Direction  EdgeDirection `json:"direction"`
	ToIndex    int           `json:"to_index"`
	Attributes `json:"attributes,omitempty"`
}

type graphJSON struct {
	Nodes []nodeJSON `json:"nodes,omitempty"`
	Edges []edgeJSON `json:"edges,omitempty"`
}

func EncodeJSON(w io.Writer, nodes Nodes) error {
	return json.NewEncoder(w).Encode(graphJSON{
		Nodes: func() []nodeJSON {
			ns := make([]nodeJSON, len(nodes))

			for i, n := range nodes {
				ns[i] = nodeJSON{
					Name:       n.Name,
					Attributes: n.Attributes,
				}
			}

			return ns
		}(),
		Edges: func() []edgeJSON {
			eix := []edgeJSON{}

			eim := map[string]struct{}{}

			for i, node := range nodes {
				for _, edge := range node.Edges {
					ei := edgeJSON{
						FromIndex: i,
						Direction: edge.Direction,
						ToIndex:   nodes.IndexOf(edge.Node),
					}
					eik := fmt.Sprintf("%#+v", ei)
					if _, ok := eim[eik]; !ok {
						eim[eik] = struct{}{}
					} else {
						continue
					}
					eix = append(eix, ei)
				}
			}

			return eix
		}(),
	})
}

func DecodeJSON(r io.Reader) (Nodes, error) {
	naej := &graphJSON{}

	err := json.NewDecoder(r).Decode(naej)
	if err != nil {
		return nil, fmt.Errorf("grap failed to decode nodes and edges JSON: %w", err)
	}

	nodes := make(Nodes, len(naej.Nodes))

	for i, naejNode := range naej.Nodes {
		nodes[i] = NewNode(naejNode.Name, naejNode.Attributes)
	}

	for _, naejEdge := range naej.Edges {
		if naejEdge.FromIndex < 0 || naejEdge.FromIndex > len(nodes) {
			continue
		}

		if naejEdge.ToIndex < 0 || naejEdge.ToIndex > len(nodes) {
			continue
		}

		var (
			name  string        = naejEdge.Name
			from  *Node         = nodes[naejEdge.FromIndex]
			to    *Node         = nodes[naejEdge.ToIndex]
			dir   EdgeDirection = naejEdge.Direction
			attrs Attributes    = naejEdge.Attributes
		)

		edge := &Edge{
			Name:       name,
			Node:       to,
			Direction:  dir,
			Attributes: attrs,
		}

		from.Edges = append(from.Edges, edge)
	}

	return nodes, nil
}
