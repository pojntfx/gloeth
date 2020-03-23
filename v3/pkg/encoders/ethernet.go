package encoders

import "net"

// Ethernet is an ethernet decoder
type Ethernet struct {
}

// NewEthernet creates a new ethernet decoder
func NewEthernet() *Ethernet {
	return &Ethernet{}
}

// GetDestMACAddress reads the destination MAC address from an ethernet frame
func (e *Ethernet) GetDestMACAddress(mac []byte) (*net.HardwareAddr, error) {
	return nil, nil
}
