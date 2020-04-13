package routing

import (
	"fmt"

	"github.com/sauerbraten/graph/v2"
)

// GetGraphFromRawData gets a graph from a multidimensional array (also see GetRawDataFromGraph)
func GetGraphFromRawData(in [][2]string) *graph.Graph {
	g := graph.New()

	adapters := DeduplicateNestedArray(in)
	nodes := GetUniqueKeys(adapters)

	fmt.Println(len(in), len(nodes), len(adapters))

	for _, node := range nodes {
		g.Add(node)
	}

	for _, adapter := range adapters {
		g.Connect(adapter[0], adapter[1], 1)
	}

	return g
}

// GetRawDataFromGraph gets a multidimensional array (also see GetGraphFromRawData)
func GetRawDataFromGraph(in *graph.Graph) [][2]string {
	nodes := in.GetAll()

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

// GetOutmostNodesFromGraph gets the outmost nodes of a graph
func GetOutmostNodesFromGraph(in *graph.Graph) []string {
	nodes := in.GetAll()

	nodeMap := make(map[string]*graph.Node)
	for _, node := range nodes {
		nodeMap[node.Key()] = node
	}

	nodesWithOneNeighbor := []string{}
	for nodeKey, node := range nodeMap {
		neighbors := node.GetNeighbors()
		if len(neighbors) == 1 {
			nodesWithOneNeighbor = append(nodesWithOneNeighbor, nodeKey)
		}
	}

	return nodesWithOneNeighbor
}
