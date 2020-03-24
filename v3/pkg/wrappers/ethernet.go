package wrappers

import (
	"net"
)

const (
	FrameSize        = 1450                         // FrameSize is the size of a non-wrapped frame
	WrappedFrameSize = 1500                         // WrappedFrameSize is the size of a wrapped frame
	HeaderSize       = WrappedFrameSize - FrameSize // HeaderSize is the size of the header
	DestSize         = 17                           // DestSize is the size of the dest address
)

// Ethernet wraps and unwraps ethernet frames
type Ethernet struct {
}

// NewEthernet creates a new ethernet wrapper
func NewEthernet() *Ethernet {
	return &Ethernet{}
}

// Wrap wraps an ethernet frame
// Format: [50]byte of header ([17]byte of dest address), [1450]byte of frame
func (e *Ethernet) Wrap(frame [FrameSize]byte, dest *net.HardwareAddr) ([WrappedFrameSize]byte, error) {
	outFrame := [WrappedFrameSize]byte{}

	outDest := [DestSize]byte{}
	copy(outDest[:], dest.String())

	outHeader := [HeaderSize]byte{}
	copy(outHeader[:DestSize], outDest[:])

	copy(outFrame[:HeaderSize], outHeader[:])
	copy(outFrame[HeaderSize:], frame[:])

	return outFrame, nil
}

// Unwrap unwraps an ethernet frame
func (e *Ethernet) Unwrap(frame [WrappedFrameSize]byte) ([FrameSize]byte, *net.HardwareAddr, error) {
	return [FrameSize]byte{}, nil, nil
}
