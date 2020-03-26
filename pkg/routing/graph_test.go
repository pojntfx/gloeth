package routing

import (
	"fmt"
	"testing"

	"github.com/sauerbraten/graph/v2"
)

var graphData = [][]string{
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

func NewGraph(graphData [][]string) *graph.Graph {
	g := graph.New()

	nodes := GetUniqueKeys(graphData)
	edges := DeduplicateNestedArray(graphData)

	fmt.Println(len(graphData), len(nodes), len(edges))

	for _, node := range nodes {
		g.Add(node)
	}

	for _, edge := range edges {
		g.Connect(edge[0], edge[1], 1)
	}

	return g
}

func DumpGraph(g *graph.Graph) [][]string {
	nodes := g.GetAll()

	nodeMap := make(map[string]*graph.Node)
	for _, node := range nodes {
		nodeMap[node.Key()] = node
	}

	var nodesWithNeighborKeys [][]string
	for nodeKey, node := range nodeMap {
		for neighbor := range node.GetNeighbors() {
			nodesWithNeighborKeys = append(nodesWithNeighborKeys, []string{nodeKey, neighbor.Key()})
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

func GetUniqueKeys(in [][]string) []string {
	outMap := make(map[string]bool)
	for _, key := range in {
		outMap[key[0]] = true
		outMap[key[1]] = true
	}

	var out []string
	for key := range outMap {
		out = append(out, key)
	}

	return out
}

func TestGetUniqueKeys(t *testing.T) {
	ins := [][][]string{
		{
			{
				"b",
				"a",
			},
			{
				"b",
				"c",
			},
			{
				"a",
				"c",
			},
			{
				"c",
				"a",
			},
			{
				"a",
				"b",
			},
			{
				"b",
				"a",
			},
		},
		{
			{
				"a",
				"b",
			},
			{
				"b",
				"c",
			},
			{
				"a",
				"c",
			},
			{
				"c",
				"a",
			},
			{
				"b",
				"a",
			},
			{
				"b",
				"a",
			},
		},
	}

	expectedOut := []string{
		"a",
		"b",
		"c",
	}

	for _, in := range ins {
		actualOut := GetUniqueKeys(in)

		actualMatchLength := 0
		expectedMatchLength := len(expectedOut)
		for _, akey := range actualOut {
			for _, ekey := range expectedOut {
				if akey == ekey {
					actualMatchLength = actualMatchLength + 1
				}
			}
		}

		t.Log(len(actualOut), len(expectedOut))
		t.Log(actualOut, expectedOut)
		t.Log(actualMatchLength, expectedMatchLength)
	}
}

func DeduplicateNestedArray(in [][]string) [][]string {
	var out [][]string

	for _, el := range in {
		match := false
		for _, nel := range out {
			if (nel[0] == el[1] && nel[1] == el[0]) || (nel[0] == el[0] && nel[1] == el[1]) {
				match = true

				break
			}
		}

		if !match {
			out = append(out, el)
		}
	}

	return out
}

func TestDeduplicateNestedArray(t *testing.T) {
	ins := [][][]string{
		{
			{
				"b",
				"a",
			},
			{
				"b",
				"c",
			},
			{
				"a",
				"c",
			},
			{
				"c",
				"a",
			},
			{
				"a",
				"b",
			},
			{
				"b",
				"a",
			},
		},
		{
			{
				"a",
				"b",
			},
			{
				"b",
				"c",
			},
			{
				"a",
				"c",
			},
			{
				"c",
				"a",
			},
			{
				"b",
				"a",
			},
			{
				"b",
				"a",
			},
		},
	}

	expectedOut := [][]string{
		{
			"b", // Note the potentially reversed value; this should not matter
			"a",
		},
		{
			"b",
			"c",
		},
		{
			"a",
			"c",
		},
	}

	for _, in := range ins {
		actualOut := DeduplicateNestedArray(in)

		actualMatchLength := 0
		expectedMatchLength := len(expectedOut)
		for _, ael := range actualOut {
			for _, eel := range expectedOut {
				if (ael[0] == eel[1] && ael[1] == eel[0]) || (ael[0] == eel[0] && ael[1] == eel[1]) {
					actualMatchLength = actualMatchLength + 1
				}
			}
		}

		t.Log(len(actualOut), len(expectedOut))
		t.Log(actualOut, expectedOut)
		t.Log(actualMatchLength, expectedMatchLength)
	}
}
