package graph_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/picatz/graph"
)

func TestEncodeDecode(t *testing.T) {
	var (
		a = graph.NewNode("a", graph.Attributes{"example": true})
		b = graph.NewNode("b", graph.Attributes{"example": "yes"})
		c = graph.NewNode("c", graph.Attributes{"example": 1})
	)

	// a → b → c

	a.AddEdgeWithDirection(b, graph.Out)
	b.AddEdgeWithDirection(c, graph.Out)

	buf := bytes.NewBuffer(nil)

	err := graph.EncodeJSON(buf, graph.Nodes{a, b, c})
	if err != nil {
		t.FailNow()
	}

	nodes, err := graph.DecodeJSON(buf)
	if err != nil {
		t.FailNow()
	}

	fmt.Println(nodes)
}
