package routing

import (
	"fmt"
	"testing"

	"github.com/sauerbraten/graph/v2"
)

var graphData = [][2]string{
	{
		"n1",
		"n4",
	},
	{
		"n2",
		"n5",
	},
	{
		"n3",
		"n6",
	},
	{
		"n6",
		"n5",
	},
	{
		"n5",
		"n4",
	},
	{
		"n8",
		"n7",
	},
	{
		"n7",
		"n4",
	},
	{
		"n7",
		"n5",
	},
	{
		"n10",
		"n9",
	},
	{
		"n9",
		"n6",
	},
	{
		"n9",
		"n7",
	},
}

func NewGraph(graphData [][2]string) *graph.Graph {
	g := graph.New()

	edges := DeduplicateNestedArray(graphData)
	nodes := GetUniqueKeys(edges)

	fmt.Println(len(graphData), len(nodes), len(edges))

	for _, node := range nodes {
		g.Add(node)
	}

	for _, edge := range edges {
		g.Connect(edge[0], edge[1], 1)
	}

	return g
}

func DumpGraph(g *graph.Graph) [][2]string {
	nodes := g.GetAll()

	nodeMap := make(map[string]*graph.Node)
	for _, node := range nodes {
		nodeMap[node.Key()] = node
	}

	var nodesWithNeighborKeys [][2]string
	for nodeKey, node := range nodeMap {
		for neighbor := range node.GetNeighbors() {
			nodesWithNeighborKeys = append(nodesWithNeighborKeys, [2]string{nodeKey, neighbor.Key()})
		}
	}

	return DeduplicateNestedArray(nodesWithNeighborKeys)
}

func TestGraph(t *testing.T) {
	g := NewGraph(graphData)

	path, err := g.ShortestPathWithHeuristic("n10", "n1", func(key, otherKey string) int {
		return 1
	})

	if err != nil {
		t.Error(err)
	}

	t.Log(path)
}

func TestDumpGraph(t *testing.T) {
	g := NewGraph(graphData)

	path, err := g.ShortestPathWithHeuristic("n10", "n1", func(key, otherKey string) int {
		return 1
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(path)

	dumpedGraph := DumpGraph(g)

	g2 := NewGraph(dumpedGraph)

	path2, err := g2.ShortestPathWithHeuristic("n10", "n1", func(key, otherKey string) int {
		return 1
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(path2)
}

func BenchmarkGraphNew(b *testing.B) {
	runs := 700

	for i := 0; i < runs; i++ {
		NewGraph(graphData)
	}
}

func BenchmarkGraph(b *testing.B) {
	runs := 700

	g := NewGraph(graphData)

	for i := 0; i < runs; i++ {
		_, err := g.ShortestPathWithHeuristic("n10", "n1", func(key, otherKey string) int {
			return 1
		})

		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGraphDump(b *testing.B) {
	runs := 700

	for i := 0; i < runs; i++ {
		g := NewGraph(graphData)

		p1, err := g.ShortestPathWithHeuristic("n10", "n1", func(key, otherKey string) int {
			return 1
		})
		if err != nil {
			b.Error(err)

			return
		}

		dumpedGraph := DumpGraph(g)

		g2 := NewGraph(dumpedGraph)

		p2, err := g2.ShortestPathWithHeuristic("n10", "n1", func(key, otherKey string) int {
			return 1
		})
		if err != nil {
			b.Error(err)

			return
		}

		b.Log(p1, p2)
	}
}
