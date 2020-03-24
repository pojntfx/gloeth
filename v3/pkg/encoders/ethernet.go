package encoders

import (
	"net"

	"github.com/pojntfx/ethernet"
	"github.com/pojntfx/gloeth/v3/pkg/encryptors"
)

// Ethernet is an ethernet decoder
type Ethernet struct {
}

// NewEthernet creates a new ethernet decoder
func NewEthernet() *Ethernet {
	return &Ethernet{}
}

// GetDestMACAddress reads the destination MAC address from an ethernet frame
func (e *Ethernet) GetDestMACAddress(frame [encryptors.PlaintextFrameSize]byte) (*net.HardwareAddr, error) {
	var ethFrame ethernet.Frame

	if err := ethFrame.UnmarshalBinary(frame[:]); err != nil {
		return nil, err
	}

	return &ethFrame.Destination, nil
}
