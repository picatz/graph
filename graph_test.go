package graph_test

import (
	"fmt"
	"testing"

	"github.com/picatz/graph"
)

// Useful stuff for drawing graphs in text format.
//
// ← ↑ → ↓ ↔ ↕ ↖ ↗ ↘ ↙
// ↰ ↱ ↲ ↳
// ⇠ ⇡ ⇢ ⇣ ⇵
// ─ ┌ ┐ └ ┘ ├ ┤ ┬ ┴ ┼ ╵ ╷ ▏▕
// ⟵ ⟶ ⟷ ⤢ ⤡ ⤾ ⤿ ⤷ ⤺ ⤻
// ╭ ╮ ╰ ╯
// ╱ ╲ ╳
// ┄
//
// ■ ▣ ▤ □ ▫ ▪
// ● ◉ ◍ ◯
// ◆ ◈ ◇
// ▲ △

func ExampleNode() {
	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}

	a.AddEdge(b)

	fmt.Println(a.Edges[0].Node.Name)
	// Output: b
}

func ExampleNode_Visit() {
	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}

	a.AddEdge(b)
	b.AddEdge(c)
	c.AddEdge(a)

	a.Visit(func(n *graph.Node) {
		fmt.Println(n.Name)
	})
	// Output: a
	// b
	// c
}

func ExampleConnectNodes() {
	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}

	graph.ConnectNodes(a, b, c)

	fmt.Println(a.PathTo(c))
	// Output: a → b → c
}

func ExampleMeshNodes() {
	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}

	//      a
	//    ⤢   ⤡
	//   b  ↔  c
	graph.MeshNodes(a, b, c)

	fmt.Println(a.PathTo(c))
	fmt.Println(c.PathTo(a))
	// Output:
	// a → c
	// c → a
}

func ExampleConnectNodes_path_to() {
	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}

	// a → b → c
	graph.ConnectNodes(a, b, c)

	path := a.PathTo(c)
	fmt.Println("a to c:", a.HasPath(c))
	fmt.Println("c to a:", c.HasPath(a))
	fmt.Println("nodes in path:", len(path))
	fmt.Println("0:", path[0].Name)
	fmt.Println("1:", path[1].Name)
	fmt.Println("2:", path[2].Name)
	// Output: a to c: true
	// c to a: false
	// nodes in path: 3
	// 0: a
	// 1: b
	// 2: c
}

func ExampleMeshNodes_path_to() {
	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}

	//      a
	//    ⤢   ⤡
	//   b  ↔  c
	graph.MeshNodes(a, b, c)

	fmt.Println("a to b:", a.HasPath(b))
	fmt.Println("b to a:", b.HasPath(a))
	fmt.Println("a to c:", a.HasPath(c))
	fmt.Println("c to a:", c.HasPath(a))
	fmt.Println("b to c:", b.HasPath(c))
	fmt.Println("c to b:", c.HasPath(b))
	// Output: a to b: true
	// b to a: true
	// a to c: true
	// c to a: true
	// b to c: true
	// c to b: true
}

func TestCyclicalGraph(t *testing.T) {
	a := &graph.Node{Name: "a"}
	b := &graph.Node{Name: "b"}
	c := &graph.Node{Name: "c"}

	a.AddEdge(b) //     a
	b.AddEdge(c) //   ↙   ↖
	c.AddEdge(a) //  b  →  c

	if !a.HasPath(a) {
		t.Fatalf("expected a to have path back to a")
	}

	if !b.HasPath(a) {
		t.Fatalf("expected b to have path back to a")
	}

	if !c.HasPath(b) {
		t.Fatalf("expected c to have path back to b")
	}

	t.Logf("a to a: %v", a.PathTo(a))
	t.Logf("b to b: %v", b.PathTo(b))
	t.Logf("c to c: %v", c.PathTo(c))
	t.Logf("a to b: %v", a.PathTo(b))
	t.Logf("a to c: %v", a.PathTo(c))
	t.Logf("b to c: %v", b.PathTo(c))
	t.Logf("b to a: %v", b.PathTo(a))
	t.Logf("c to a: %v", c.PathTo(a))
	t.Logf("c to b: %v", c.PathTo(b))

	t.Logf("a has cycles: %v", a.HasCycles())
	t.Logf("b has cycles: %v", b.HasCycles())
	t.Logf("c has cycles: %v", c.HasCycles())
}

func TestFindBridges(t *testing.T) {
	tests := []struct {
		Name    string
		Root    *graph.Node
		Bridges map[string]bool
	}{
		{
			Name: "wolfram",
			Bridges: map[string]bool{
				"e → b": true,
				"f → b": true,
				"d → a": true,
			},
			Root: func() *graph.Node {
				a := &graph.Node{Name: "a"}
				b := &graph.Node{Name: "b"}
				c := &graph.Node{Name: "c"}
				d := &graph.Node{Name: "d"}
				e := &graph.Node{Name: "e"}
				f := &graph.Node{Name: "f"}

				//        a ← d
				//      ↙   ↖
				// e → b  →  c
				//     ↑
				//     f

				a.AddEdge(b)
				b.AddEdge(c)
				c.AddEdge(a)
				d.AddEdge(a)
				e.AddEdge(b)
				f.AddEdge(b)

				return a
			}(),
		},
		{
			Name: "TIE fighter single direction",
			Bridges: map[string]bool{
				"c → d": true,
			},
			Root: func() *graph.Node {
				a := &graph.Node{Name: "a"}
				b := &graph.Node{Name: "b"}
				c := &graph.Node{Name: "c"}
				d := &graph.Node{Name: "d"}
				e := &graph.Node{Name: "e"}
				f := &graph.Node{Name: "f"}

				// a           e
				// ↑ ⤡       ⤢ ↑
				// |   c → d   |
				// ↓ ⤢       ⤡ ↓
				// b           f

				a.AddLink(b)
				c.AddLink(a)
				c.AddLink(b)
				c.AddEdge(d)
				d.AddLink(e)
				d.AddLink(f)
				f.AddLink(e)
				return a
			}(),
		},
		{
			Name:    "TIE fighter bi-directional",
			Bridges: map[string]bool{
				// c ↔ d is a bi-directional relationship
				//
				// So, it's really two edges, a "bridge pair".
				//
				// c → d
				// d → c
				//
				// If one of those edges is removed, the
				// bridge is still maintained with the other.
			},
			Root: func() *graph.Node {
				a := &graph.Node{Name: "a"}
				b := &graph.Node{Name: "b"}
				c := &graph.Node{Name: "c"}
				d := &graph.Node{Name: "d"}
				e := &graph.Node{Name: "e"}
				f := &graph.Node{Name: "f"}

				// a           e
				// ↑ ⤡       ⤢ ↑
				// |   c ↔ d   |
				// ↓ ⤢       ⤡ ↓
				// b           f

				a.AddLink(b)
				c.AddLink(a)
				c.AddLink(b)
				c.AddLink(d)
				d.AddLink(e)
				d.AddLink(f)
				f.AddLink(e)
				return a
			}(),
		},
		{
			Name: "tree",
			Bridges: map[string]bool{
				"a → c": true,
				"a → b": true,
				"b → d": true,
				"b → e": true,
				"e → f": true,
			},
			Root: func() *graph.Node {
				a := &graph.Node{Name: "a"}
				b := &graph.Node{Name: "b"}
				c := &graph.Node{Name: "c"}
				d := &graph.Node{Name: "d"}
				e := &graph.Node{Name: "e"}
				f := &graph.Node{Name: "f"}

				//       a
				//     ↙   ↘
				//    b     c
				//  ↙   ↘
				// d     e
				//       ↓
				//       f

				a.AddEdge(b)
				a.AddEdge(c)
				b.AddEdge(d)
				b.AddEdge(e)
				e.AddEdge(f)

				return a
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			bridges := graph.FindBridges(test.Root)

			for _, bridge := range bridges {
				_, ok := test.Bridges[bridge.String()]
				if !ok {
					t.Logf("unexpected bridge found: %v", bridge)
					t.Fail()
				}
			}

			if len(bridges) != len(test.Bridges) {
				t.Logf("unexpected number of bridges found: expected: %d, got: %d", len(test.Bridges), len(bridges))
				t.Fail()
			}
		})
	}
}
