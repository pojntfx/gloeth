package wrappers

import (
	"net"
)

const (
	WrappedFrameSize   = 1780                                        // WrappedFrameSize is the size of a wrapped frame
	HeaderSize         = 250                                         // HeaderSize is the size of the header
	HeaderDestSize     = 17                                          // DestSize is the size of the dest address
	HeaderSrcSize      = HeaderDestSize                              // SrcSize is the size of the dest address
	HopSize            = HeaderDestSize                              // HopSize is the size of a hop address
	HopsSize           = HeaderSize - HeaderDestSize - HeaderSrcSize // HopsSize is the size of the hops
	HopsCount          = 12                                          // HopsCount is the maximum amount of hops
	EncryptedFrameSize = WrappedFrameSize - HeaderSize               // EncryptedFrameSize is the size of a encrypted, non-wrapped frame
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
func (e *Ethernet) Wrap(dest, src *net.HardwareAddr, hops [HopsCount]*net.HardwareAddr, frame [EncryptedFrameSize]byte) ([WrappedFrameSize]byte, error) {
	outFrame := [WrappedFrameSize]byte{}

	outDest := [HeaderDestSize]byte{}
	copy(outDest[:], dest.String())

	outSrc := [HeaderSrcSize]byte{}
	copy(outSrc[:], src.String())

	outHops := [HopsSize]byte{}
	for i, hop := range hops {
		if hop == nil {
			continue
		}

		outHop := [HopSize]byte{}
		copy(outHop[:], hop.String())

		copy(outHops[i*HopSize:(i+1)*HopSize], outHop[:])
	}

	outHeader := [HeaderSize]byte{}
	copy(outHeader[:HeaderDestSize], outDest[:])
	copy(outHeader[HeaderDestSize:HeaderDestSize+HeaderSrcSize], outSrc[:])
	copy(outHeader[HeaderDestSize+HeaderSrcSize:HeaderDestSize+HeaderSrcSize+HopsSize], outHops[:])

	copy(outFrame[:HeaderSize], outHeader[:])
	copy(outFrame[HeaderSize:], frame[:])

	return outFrame, nil
}

// Unwrap unwraps an ethernet frame
func (e *Ethernet) Unwrap(frame [WrappedFrameSize]byte) (*net.HardwareAddr, *net.HardwareAddr, [HopsCount]*net.HardwareAddr, [EncryptedFrameSize]byte, error) {
	outHeader := frame[:HeaderSize]

	outDest, err := net.ParseMAC(string(outHeader[:HeaderDestSize]))
	if err != nil {
		return nil, nil, [HopsCount]*net.HardwareAddr{}, [EncryptedFrameSize]byte{}, err
	}

	outSrc, err := net.ParseMAC(string(outHeader[HeaderDestSize : HeaderDestSize+HeaderSrcSize]))
	if err != nil {
		return &outDest, nil, [HopsCount]*net.HardwareAddr{}, [EncryptedFrameSize]byte{}, err
	}

	outHops := [HopsCount]*net.HardwareAddr{}
	for i := 0; i < HopsCount; i++ {
		inHop := outHeader[HeaderDestSize+HeaderSrcSize+(i*HopSize) : HeaderDestSize+HeaderSrcSize+((i+1)*HopSize)]

		isEmpty := true
		for i := 0; i < HopsCount; i++ {
			if inHop[i] != 0 {
				isEmpty = false

				break
			}
		}

		if isEmpty {
			continue
		}

		outHop, err := net.ParseMAC(string(inHop))
		if err != nil {
			return nil, nil, [HopsCount]*net.HardwareAddr{}, [EncryptedFrameSize]byte{}, err
		}

		outHops[i] = &outHop
	}

	outFrame := [EncryptedFrameSize]byte{}
	copy(outFrame[:], frame[HeaderSize:])

	return &outDest, &outSrc, outHops, outFrame, nil
}

// GetShiftedHops returns the hops for the next switcher
func (e *Ethernet) GetShiftedHops(hops [HopsCount]*net.HardwareAddr) [HopsCount]*net.HardwareAddr {
	outHops := [HopsCount]*net.HardwareAddr{}

	for i, hop := range hops {
		if i != 0 {
			outHops[i-1] = hop
		}
	}

	return outHops
}

// GetHopsEmpty returns true if there are no more hops
func (e *Ethernet) GetHopsEmpty(hops [HopsCount]*net.HardwareAddr) bool {
	isEmpty := true

	for _, hop := range hops {
		if hop != nil {
			isEmpty = false

			break
		}
	}

	return isEmpty
}
