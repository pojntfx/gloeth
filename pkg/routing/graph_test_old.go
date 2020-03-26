package routing

import (
	"testing"
)

var graphData = [][2]string{}

func TestGraph(t *testing.T) {
	g := GetGraphFromRawData(graphData)

	path, err := g.ShortestPathWithHeuristic("n10", "n1", func(key, otherKey string) int {
		return 1
	})

	if err != nil {
		t.Error(err)
	}

	t.Log(path)
}

func TestDumpGraph(t *testing.T) {
	g := GetGraphFromRawData(graphData)

	path, err := g.ShortestPathWithHeuristic("n10", "n1", func(key, otherKey string) int {
		return 1
	})
	if err != nil {
		t.Error(err)
	}

	t.Log(path)

	dumpedGraph := GetRawDataFromGraph(g)

	g2 := GetGraphFromRawData(dumpedGraph)

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
		GetGraphFromRawData(graphData)
	}
}

func BenchmarkGraph(b *testing.B) {
	runs := 700

	g := GetGraphFromRawData(graphData)

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
		g := GetGraphFromRawData(graphData)

		p1, err := g.ShortestPathWithHeuristic("n10", "n1", func(key, otherKey string) int {
			return 1
		})
		if err != nil {
			b.Error(err)

			return
		}

		dumpedGraph := GetRawDataFromGraph(g)

		g2 := GetGraphFromRawData(dumpedGraph)

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
