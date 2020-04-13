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
