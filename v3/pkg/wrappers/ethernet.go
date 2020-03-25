package wrappers

import (
	"net"
)

const (
	EncryptedFrameSize = 1450                                  // EncryptedFrameSize is the size of a encrypted, non-wrapped frame
	WrappedFrameSize   = 1500                                  // WrappedFrameSize is the size of a wrapped frame
	HeaderSize         = WrappedFrameSize - EncryptedFrameSize // HeaderSize is the size of the header
	DestSize           = 17                                    // DestSize is the size of the dest address
	SrcSize            = 17                                    // SrcSize is the size of the dest address
)

// Ethernet wraps and unwraps ethernet frames
type Ethernet struct {
}

// NewEthernet creates a new ethernet wrapper
func NewEthernet() *Ethernet {
	return &Ethernet{}
}

// Wrap wraps an ethernet frame
// Format: [HeaderSize]byte of header ([DestSize]byte of dest address, [SrcSize]byte of src address), [EncryptedFrameSize]byte of frame
func (e *Ethernet) Wrap(dest, src *net.HardwareAddr, frame [EncryptedFrameSize]byte) ([WrappedFrameSize]byte, error) {
	outFrame := [WrappedFrameSize]byte{}

	outDest := [DestSize]byte{}
	copy(outDest[:], dest.String())

	outSrc := [SrcSize]byte{}
	copy(outSrc[:], src.String())

	outHeader := [HeaderSize]byte{}
	copy(outHeader[:DestSize], outDest[:])
	copy(outHeader[DestSize:DestSize+SrcSize], outSrc[:])

	copy(outFrame[:HeaderSize], outHeader[:])
	copy(outFrame[HeaderSize:], frame[:])

	return outFrame, nil
}

// Unwrap unwraps an ethernet frame
func (e *Ethernet) Unwrap(frame [WrappedFrameSize]byte) (*net.HardwareAddr, *net.HardwareAddr, [EncryptedFrameSize]byte, error) {
	outHeader := frame[:HeaderSize]

	outDest, err := net.ParseMAC(string(outHeader[:DestSize]))
	if err != nil {
		return nil, nil, [EncryptedFrameSize]byte{}, err
	}

	outSrc, err := net.ParseMAC(string(outHeader[DestSize : DestSize+SrcSize]))
	if err != nil {
		return &outDest, nil, [EncryptedFrameSize]byte{}, err
	}

	outFrame := [EncryptedFrameSize]byte{}
	copy(outFrame[:], frame[HeaderSize:])

	return &outDest, &outSrc, outFrame, nil
}
