package devices

import (
	"github.com/pojntfx/gloeth/v3/pkg/encryptors"
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

const (
	MTU = encryptors.PlaintextFrameSize - 14 // MTU is the MTU, which is the plaintext frame size - ethernet header (14 bytes)
)

// TAP is a TAP device
type TAP struct {
	readChan chan [encryptors.PlaintextFrameSize]byte
	mtu      uint
	name     string
	dev      *water.Interface
}

// NewTAP creates a new TAP device
func NewTAP(readChan chan [encryptors.PlaintextFrameSize]byte, mtu uint, name string) *TAP {
	return &TAP{readChan, mtu, name, nil}
}

// Open opens the TAP device
func (t *TAP) Open() error {
	dev, err := water.New(water.Config{
		DeviceType: water.TAP,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: t.name,
		},
	})
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(t.name)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetMTU(link, int(t.mtu)); err != nil {
		return err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	t.dev = dev

	return nil
}

// Close closes the TAP device
func (t *TAP) Close() error {
	return t.dev.Close()
}

// Read reads from the TAP device
func (t *TAP) Read() error {
	for {
		readFrame := [encryptors.PlaintextFrameSize]byte{}

		_, err := t.dev.Read(readFrame[:])
		if err != nil {
			return err
		}

		t.readChan <- readFrame
	}
}

// Write writes from the TAP device
func (t *TAP) Write(frame [encryptors.PlaintextFrameSize]byte) error {
	_, err := t.dev.Write(frame[:])

	return err
}
