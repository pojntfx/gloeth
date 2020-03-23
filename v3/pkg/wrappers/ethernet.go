package wrappers

import "net"

// Ethernet wraps and unwraps ethernet frames
type Ethernet struct {
}

// NewEthernet creates a new ethernet wrapper
func NewEthernet() *Ethernet {
	return &Ethernet{}
}

// Wrap wraps an ethernet frame
func (e *Ethernet) Wrap(frame []byte, dest *net.HardwareAddr) ([]byte, error) {
	return nil, nil
}

// Unwrap unwraps an ethernet frame
func (e *Ethernet) Unwrap(frame []byte) ([]byte, *net.HardwareAddr, error) {
	return nil, nil, nil
}
