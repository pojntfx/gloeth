package routing

import (
	"errors"
	"net"

	"github.com/sauerbraten/graph/v2"
)

// RoutingTable manages routing information
type RoutingTable struct {
	graph *graph.Graph
}

// NewRoutingTable creates a new routing table
func NewRoutingTable() *RoutingTable {
	return &RoutingTable{
		graph.New(),
	}
}

// Register adds and connects two nodes
func (r *RoutingTable) Register(mac1, mac2 *net.HardwareAddr) error {
	mac1Key, mac2Key := mac1.String(), mac2.String()

	r.graph.Add(mac1Key)
	r.graph.Add(mac2Key)

	ok := r.graph.Connect(mac1Key, mac2Key, 1)

	if !ok {
		return errors.New("could not connect nodes, either both are the same or the keys are invalid")
	}

	return nil
}

// GetHops returns the hops between a switcher and an adapter
func (r *RoutingTable) GetHops(switcherMAC, adapterMAC *net.HardwareAddr) ([]*net.HardwareAddr, error) {
	fullPath, err := r.graph.ShortestPathWithHeuristic(switcherMAC.String(), adapterMAC.String(), func(key, otherKey string) int {
		return 1
	})

	if err != nil {
		return []*net.HardwareAddr{}, err
	}

	hops := make([]*net.HardwareAddr, len(fullPath))
	for i, rawHop := range fullPath {
		hop, err := net.ParseMAC(rawHop)
		if err != nil {
			return []*net.HardwareAddr{}, err
		}

		hops[(len(fullPath)-1)-i] = &hop
	}

	if len(hops) <= 2 {
		return []*net.HardwareAddr{}, nil
	}

	return hops[1 : len(fullPath)-1], nil
}

// Marshal returns the routing table as raw data
func (r *RoutingTable) Marshal() [][2]string {
	return GetRawDataFromGraph(r.graph)
}

// Unmarshal creates the routing table from raw data
func (r *RoutingTable) Unmarshal(rawData [][2]string) {
	r.graph = GetGraphFromRawData(rawData)
}
