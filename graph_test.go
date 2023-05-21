package graph_test

import (
	"fmt"
	"reflect"
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

	graph.AddEdges(
		graph.AddEdge{From: a, To: b},
		graph.AddEdge{From: b, To: c},
		graph.AddEdge{From: c, To: a},
	)

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

func TestAddEdgeWithDirection(t *testing.T) {
	a := &graph.Node{}
	b := &graph.Node{}

	a.AddEdgeWithDirection(b, graph.Out) // a  →  b

	if !a.HasPath(b) {
		t.Fatalf("expected a to have path to b")
	}

	if b.HasPath(a) {
		t.Fatalf("did not expect b to have path back to a")
	}
}

func TestDirection(t *testing.T) {
	tests := []struct {
		Name      string
		Direction graph.EdgeDirection
		String    string
	}{
		{
			Name:      "unknown",
			Direction: graph.Unknown,
			String:    "┄",
		},
		{
			Name:      "none",
			Direction: graph.None,
			String:    "-",
		},
		{
			Name:      "in",
			Direction: graph.In,
			String:    "←",
		},
		{
			Name:      "out",
			Direction: graph.Out,
			String:    "→",
		},
		{
			Name:      "both",
			Direction: graph.Both,
			String:    "↔",
		},
		{
			Name:      "anything else",
			Direction: graph.EdgeDirection(100),
			String:    "┄",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.Direction.String() != test.String {
				t.Fatalf("expected: %q, got: %q", test.String, test.Direction)
			}
		})
	}
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
	only := []string{
		// "simple dangling edge",
		// "multiple dangling edges next to a cycle",
		// "dangling edge next to a cycle",
		// "nested dangling edge next to a cycle",
		// "TIE fighter (barbell) single direction",
		// "TIE fighter (barbell) bi-directional",
		// "tree",
	}

	tests := []struct {
		Name    string
		Root    *graph.Node
		Bridges map[string]bool
	}{
		{
			Name: "simple dangling edge",
			Bridges: map[string]bool{
				"a → b": true,
			},
			Root: func() *graph.Node {
				a := &graph.Node{Name: "a"}
				b := &graph.Node{Name: "b"}

				// a → b

				a.AddEdge(b)

				return a
			}(),
		},
		{
			Name: "dangling edge next to a cycle",
			Bridges: map[string]bool{
				"a → b": true,
			},
			Root: func() *graph.Node {
				a := &graph.Node{Name: "a"}
				b := &graph.Node{Name: "b"}
				c := &graph.Node{Name: "c"}
				d := &graph.Node{Name: "d"}

				//        c
				//      ↗   ↘
				// a → b  ←  d
				//

				a.AddEdge(b)
				b.AddEdge(c)
				c.AddEdge(d)
				d.AddEdge(b)

				return a
			}(),
		},
		{
			Name: "nested dangling edge next to a cycle",
			Bridges: map[string]bool{
				"a → b": true,
				"b → c": true,
			},
			Root: func() *graph.Node {
				a := &graph.Node{Name: "a"}
				b := &graph.Node{Name: "b"}
				c := &graph.Node{Name: "c"}
				d := &graph.Node{Name: "d"}
				e := &graph.Node{Name: "e"}

				//            d
				//          ↗   ↘
				// a → b → c  ←  e
				//

				a.AddEdge(b)
				b.AddEdge(c)
				c.AddEdge(d)
				d.AddEdge(e)
				e.AddEdge(c)

				return a
			}(),
		},
		{
			Name: "multiple dangling edges next to a cycle",
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
			Name: "TIE fighter (barbell) single direction",
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
				c.AddEdge(d) // this is the bridge
				d.AddLink(e)
				d.AddLink(f)
				f.AddLink(e)
				return a
			}(),
		},
		{
			Name:    "TIE fighter (barbell) bi-directional",
			Bridges: map[string]bool{
				// c ↔ d is a bi-directional relationship
				//
				// It's really two edges. If one of those edges
				// is removed, the bridge is still maintained
				// with the other.
				//
				// "c → d": false,
				// "d → c": false,
				//
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
		// To optionally debug only certain tests, useful
		// for interactive debugging without needing to skip
		// graphs you're not interested in.
		if len(only) > 0 {
			var run bool
			for _, onlyTest := range only {
				if onlyTest == test.Name {
					run = true
					break
				}
			}
			if !run {
				continue
			}
		}
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
				for _, bridge := range bridges {
					t.Logf("\t%s", bridge)
				}
				t.Fail()
			}
		})
	}
}

func TestFindAdjacentTo(t *testing.T) {
	tests := []struct {
		Name       string
		Case       func(*testing.T)
		Unexpected bool
	}{
		{
			Name:       "simple",
			Unexpected: true,
			Case: func(t *testing.T) {
				var (
					a = &graph.Node{Name: "a"}
					b = &graph.Node{Name: "b"}
					c = &graph.Node{Name: "c"}
				)

				// a  →  b  →  c

				a.AddEdge(b)
				b.AddEdge(c)

				if !a.Edges.AdjacentTo(b) {
					t.Fail()
				}

				if !b.Edges.AdjacentTo(c) {
					t.Fail()
				}

				if a.Edges.AdjacentTo(c) {
					t.Fail()
				}
			},
		},
		{
			Name:       "complex",
			Unexpected: true,
			Case: func(t *testing.T) {
				var (
					a = &graph.Node{Name: "a"}
					b = &graph.Node{Name: "b"}
					c = &graph.Node{Name: "c"}
					d = &graph.Node{Name: "d"}
					e = &graph.Node{Name: "e"}
				)

				//           b
				//         ↙   ↖
				//       c       a
				//     ↙   ↘   ↗
				//    e  →   d

				a.AddEdge(b)
				b.AddEdge(c)
				c.AddEdge(d)
				d.AddEdge(a)
				c.AddEdge(e)
				e.AddEdge(d)

				if !a.Edges.AdjacentTo(b) {
					t.Fail()
				}

				if !b.Edges.AdjacentTo(c) {
					t.Fail()
				}

				if a.Edges.AdjacentTo(c) {
					t.Fail()
				}

				if !c.Edges.AdjacentTo(d) {
					t.Fail()
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, test.Case)
	}
}

func TestFindCliques(t *testing.T) {
	var (
		a = &graph.Node{Name: "a"}
		b = &graph.Node{Name: "b"}
		c = &graph.Node{Name: "c"}
		d = &graph.Node{Name: "d"}
		e = &graph.Node{Name: "e"}
		f = &graph.Node{Name: "f"}
		g = &graph.Node{Name: "g"}
		h = &graph.Node{Name: "h"}
		i = &graph.Node{Name: "i"}
		j = &graph.Node{Name: "j"}
		k = &graph.Node{Name: "k"}
		l = &graph.Node{Name: "l"}
		m = &graph.Node{Name: "m"}
	)

	//       a
	//     ↙   ↘
	//    b  →  c
	//            ↘
	//   i  ←  h    d → e
	//    ↘   ↗   ↙
	//      g ← f
	//    ↙ ↑ ↘
	//  j   m ← l
	//    ↘ ↑ ↗
	//      k
	//
	// Cliques (of 3 or more)
	//
	// 1. a, b, c
	// 2. g, h, i
	// 3. g, m, l
	// 4. k, m, l
	//

	a.AddEdge(b)
	a.AddEdge(c)
	b.AddEdge(c)
	c.AddEdge(d)
	d.AddEdge(e)
	d.AddEdge(f)
	f.AddEdge(g)
	g.AddEdge(h)
	h.AddEdge(i)
	i.AddEdge(g)
	g.AddEdge(j)
	j.AddEdge(k)
	k.AddEdge(l)
	k.AddEdge(m)
	l.AddEdge(m)
	g.AddEdge(l)
	m.AddEdge(g)

	cliques := graph.FindCliques(a, 3)

	t.Logf("found %d cliques", len(cliques))
	for _, clique := range cliques {
		t.Logf("clique: %v", clique)
	}
}

func TestFindCliques_2(t *testing.T) {
	var (
		a = &graph.Node{Name: "a"}
		b = &graph.Node{Name: "b"}
		c = &graph.Node{Name: "c"}
		d = &graph.Node{Name: "d"}
		e = &graph.Node{Name: "e"}
	)

	//           b
	//         ↙   ↖
	//       c       a
	//     ↙   ↘   ↗
	//    e  →   d
	//

	a.AddEdge(b)
	b.AddEdge(c)
	c.AddEdge(d)
	d.AddEdge(a)
	c.AddEdge(e)
	e.AddEdge(d)

	cliques := graph.FindCliques(a, 3)

	if len(cliques) != 1 {
		t.Fail()
	}

	if len(cliques[0].Nodes()) != 3 {
		t.Fail()
	}

	t.Logf("found %d cliques", len(cliques))
	for _, clique := range cliques {
		t.Logf("clique: %v", clique)
	}
}

func TestAttributes(t *testing.T) {
	attrs := graph.Attributes{
		"hello":   "world",
		"enabled": true,
		"size":    100,
		"rate":    1.0,
	}

	err := graph.UseAttribute(attrs, "hello", func(v string) {
		if v != "world" {
			t.Fail()
		}
	})
	if err != nil {
		t.Fail()
	}

	err = graph.UseAttribute(attrs, "enabled", func(v bool) {
		if v != true {
			t.Fail()
		}
	})
	if err != nil {
		t.Fail()
	}

	err = graph.UseAttribute(attrs, "size", func(v int) {
		if v != 100 {
			t.Fail()
		}
	})
	if err != nil {
		t.Fail()
	}

	err = graph.UseAttribute(attrs, "rate", func(v float64) {
		if v != 1.0 {
			t.Fail()
		}
	})
	if err != nil {
		t.Fail()
	}
}

func TestIsBipartite_false(t *testing.T) {
	var (
		a = graph.NewNode("a", nil)
		b = graph.NewNode("b", nil)
		c = graph.NewNode("c", nil)
		d = graph.NewNode("d", nil)
		e = graph.NewNode("e", nil)
	)

	//           b
	//         ↙   ↖
	//       c       a
	//     ↙   ↘   ↗
	//    e  →   d
	//

	a.AddEdge(b)
	b.AddEdge(c)
	c.AddEdge(d)
	d.AddEdge(a)
	c.AddEdge(e)
	e.AddEdge(d)

	g := graph.New("test", graph.WithNodes(graph.NewNodes(
		a, b, c, d, e,
	)))

	if g.IsBipartite() {
		t.Fail()
	}
}

func TestIsBipartite_true(t *testing.T) {
	var (
		a = graph.NewNode("a", nil)
		b = graph.NewNode("b", nil)
		c = graph.NewNode("c", nil)
		d = graph.NewNode("d", nil)
		e = graph.NewNode("e", nil)
	)

	//  a   b   c
	//   ↘ ↙ ↘ ↙
	//    d   e

	a.AddEdge(d)
	b.AddEdge(d)
	b.AddEdge(e)
	c.AddEdge(e)

	g := graph.New("test", graph.WithNodes(graph.NewNodes(
		a, b, c, d, e,
	)))

	if !g.IsBipartite() {
		t.Fail()
	}
}

func TestIsMultipartite_2_true(t *testing.T) {
	var (
		a = graph.NewNode("a", nil)
		b = graph.NewNode("b", nil)
		c = graph.NewNode("c", nil)
		d = graph.NewNode("d", nil)
		e = graph.NewNode("e", nil)
	)

	//  a   b   c
	//   ↘ ↙ ↘ ↙
	//    d   e

	a.AddEdge(d)
	b.AddEdge(d)
	b.AddEdge(e)
	c.AddEdge(e)

	g := graph.New("test", graph.WithNodes(graph.NewNodes(
		a, b, c, d, e,
	)))

	if !g.IsMultipartite(2) {
		t.Fail()
	}
}

func TestIsMultipartite_2_false(t *testing.T) {
	var (
		a = graph.NewNode("a", nil)
		b = graph.NewNode("b", nil)
		c = graph.NewNode("c", nil)
		d = graph.NewNode("d", nil)
		e = graph.NewNode("e", nil)
	)

	//  a   b   c
	//   ↘ ↙ ↘ ↙
	//    d → e

	a.AddEdge(d)
	b.AddEdge(d)
	b.AddEdge(e)
	c.AddEdge(e)
	d.AddEdge(e)

	g := graph.New("test", graph.WithNodes(graph.NewNodes(
		a, b, c, d, e,
	)))

	if g.IsMultipartite(2) {
		t.Fail()
	}
}

func TestInstance_DFS(t *testing.T) {
	// Create a new graph instance.
	inst := graph.New("test")

	// Add some nodes to the graph.
	nodeA := graph.NewNode("a", nil)
	nodeB := graph.NewNode("b", nil)
	nodeC := graph.NewNode("c", nil)
	nodeD := graph.NewNode("d", nil)
	nodeE := graph.NewNode("e", nil)

	inst.AddNodes(
		nodeA,
		nodeB,
		nodeC,
		nodeD,
		nodeE,
	)

	// Add some edges to the graph.
	//
	//  ┌───────────────┐
	//  ↓               │
	//  a → b → c → d → e
	//
	nodeA.AddEdge(nodeB)
	nodeB.AddEdge(nodeC)
	nodeC.AddEdge(nodeD)
	nodeD.AddEdge(nodeE)
	nodeE.AddEdge(nodeA)

	// Create a slice to store the visited nodes.
	visited := graph.Nodes{}

	// Define a function to visit each node.
	fn := func(node *graph.Node) {
		visited = append(visited, node)
	}

	// Perform DFS on the graph.
	inst.DFS(fn)

	// Check that the visited nodes are in the correct order.
	expected := graph.Nodes{
		nodeA,
		nodeB,
		nodeC,
		nodeD,
		nodeE,
	}

	if !reflect.DeepEqual(visited, expected) {
		t.Fail()
	}
}

func TestInstance_BFS(t *testing.T) {
	// Create a new graph instance.
	inst := graph.New("test")

	// Add some nodes to the graph.
	nodeA := graph.NewNode("a", nil)
	nodeB := graph.NewNode("b", nil)
	nodeC := graph.NewNode("c", nil)
	nodeD := graph.NewNode("d", nil)
	nodeE := graph.NewNode("e", nil)
	nodeF := graph.NewNode("f", nil)
	nodeG := graph.NewNode("g", nil)
	nodeH := graph.NewNode("h", nil)

	inst.AddNodes(
		nodeA,
		nodeB,
		nodeC,
		nodeD,
		nodeE,
		nodeF,
		nodeG,
		nodeH,
	)

	// Add some edges to the graph.
	//
	//     c           h
	//     ↑           ↑
	// a → b → d → f → g
	//     ↓   |
	//     e ←─┘
	//
	nodeA.AddEdge(nodeB)
	nodeB.AddEdge(nodeC)
	nodeB.AddEdge(nodeD)
	nodeB.AddEdge(nodeE)
	nodeD.AddEdge(nodeE)
	nodeD.AddEdge(nodeF)
	nodeF.AddEdge(nodeG)
	nodeG.AddEdge(nodeH)

	// Create a slice to store the visited nodes.
	visited := graph.Nodes{}

	// Perform BFS on the graph.
	inst.BFS(func(node *graph.Node) {
		visited = append(visited, node)
	})

	// Check that the visited nodes are in the correct order.
	expected := graph.Nodes{
		nodeA,
		nodeB,
		nodeC,
		nodeD,
		nodeE,
		nodeF,
		nodeG,
		nodeH,
	}

	if !reflect.DeepEqual(visited, expected) {
		t.Errorf("visited nodes = %v, expected %v", visited, expected)
	}
}
