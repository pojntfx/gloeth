package encoders

import (
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/pojntfx/gloeth/pkg/encryptors"
)

// Ethernet is an ethernet decoder
type Ethernet struct {
}

// NewEthernet creates a new ethernet decoder
func NewEthernet() *Ethernet {
	return &Ethernet{}
}

// GetMACAddresses reads the destination and source MAC addresses from an ethernet frame
func (e *Ethernet) GetMACAddresses(frame [encryptors.PlaintextFrameSize]byte) (*net.HardwareAddr, *net.HardwareAddr, error) {
	var ethFrame ethernet.Frame

	if err := ethFrame.UnmarshalBinary(frame[:]); err != nil {
		return nil, nil, err
	}

	return &ethFrame.Destination, &ethFrame.Source, nil
}
