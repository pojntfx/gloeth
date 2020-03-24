package wrappers

import (
	"fmt"
	"net"
	"testing"

	gm "github.com/cseeger-epages/mac-gen-go"
)

func getDest() (net.HardwareAddr, error) {
	prefix := gm.GenerateRandomLocalMacPrefix(false)
	suffix, err := gm.CalculateNICSufix(net.ParseIP("10.0.0.1"))
	if err != nil {
		return nil, err
	}

	rawDest := fmt.Sprintf("%v:%v", prefix, suffix)

	return net.ParseMAC(rawDest)
}

func getFrame() [FrameSize]byte {
	return [FrameSize]byte{1}
}

func getWrappedFrame(dest net.HardwareAddr, frame [FrameSize]byte) [WrappedFrameSize]byte {
	outFrame := [WrappedFrameSize]byte{}

	outDest := [DestSize]byte{}
	copy(outDest[:], dest.String())

	outHeader := [HeaderSize]byte{}
	copy(outHeader[:DestSize], outDest[:])

	copy(outFrame[:HeaderSize], outHeader[:])
	copy(outFrame[HeaderSize:], frame[:])

	return outFrame
}

func TestNewEthernet(t *testing.T) {
	e := NewEthernet()

	if e == nil {
		t.Error("Ethernet is nil")
	}
}

func TestWrap(t *testing.T) {
	expectedFrame := getFrame()
	expectedDest, err := getDest()
	if err != nil {
		t.Error(err)
	}

	e := NewEthernet()

	wrappedFrame, err := e.Wrap(&expectedDest, expectedFrame)
	if err != nil {
		t.Error(err)
	}

	actualHeader := wrappedFrame[:HeaderSize]
	actualDest, err := net.ParseMAC(string(actualHeader[:DestSize]))
	if err != nil {
		t.Error(err)
	}

	if actualDest.String() != expectedDest.String() {
		t.Error("Dest not wrapped correctly")
	}

	actualFrame := [FrameSize]byte{}
	copy(actualFrame[:], wrappedFrame[HeaderSize:])
	if actualFrame != expectedFrame {
		t.Error("Frame not wrapped correctly")
	}
}

// TODO: Add test for faulty MAC address
func TestUnwrap(t *testing.T) {
	expectedFrame := getFrame()
	expectedDest, err := getDest()
	if err != nil {
		t.Error(err)
	}
	wrappedFrame := getWrappedFrame(expectedDest, expectedFrame)

	e := NewEthernet()

	actualDest, actualFrame, err := e.Unwrap(wrappedFrame)
	if err != nil {
		t.Error(err)
	}

	if actualDest.String() != expectedDest.String() {
		t.Error("Dest not unwrapped correctly")
	}

	if actualFrame != expectedFrame {
		t.Error("Frame not unwrapped correctly")
	}
}
