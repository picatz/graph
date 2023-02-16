package graph

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func EncodeDOT(w io.Writer, nodes Nodes) error {
	var err error

	bw := bufio.NewWriter(w)

	bw.WriteString("digraph {\n")

	for _, node := range nodes {
		if len(node.Edges.Out()) > 0 {
			_, err = bw.WriteString(
				fmt.Sprintf(
					"\t%q -> { %s }\n",
					node.Name,
					func() string {
						var names []string
						for _, edge := range node.Edges.Out() {
							names = append(names, fmt.Sprintf("%q", edge.Node.Name))
						}
						return strings.Join(names, " ")
					}(),
				),
			)
			if err != nil {
				return fmt.Errorf("graph failed to encode DOT: %w", err)
			}
		}
	}

	bw.WriteString("}\n")

	err = bw.Flush()
	if err != nil {
		return fmt.Errorf("graph failed to encode DOT: %w", err)
	}
	return nil
}

func DecodeDOT(r io.Reader) (Nodes, error) {
	return nil, fmt.Errorf("graph decode DOT not implemented yet")
}
