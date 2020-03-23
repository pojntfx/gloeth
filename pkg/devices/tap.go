package devices

import (
	"github.com/pojntfx/gloeth/pkg/constants"
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

// TAPDevice is a TAP device
type TAPDevice struct {
	mtu        int
	framesChan chan []byte
	iface      *water.Interface
	name       string
}

// NewTAPDevice creates a new TAP device
func NewTAPDevice(mtu int, name string, framesChan chan []byte) *TAPDevice {
	return &TAPDevice{
		mtu:        mtu,
		name:       name,
		framesChan: framesChan,
	}
}

// Open opens the TAP device
func (t *TAPDevice) Open() error {
	dev, err := water.New(water.Config{DeviceType: water.TAP, PlatformSpecificParams: water.PlatformSpecificParams{Name: t.name}})
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(t.name)
	if err != nil {
		return err
	}
	if err := netlink.LinkSetMTU(link, t.mtu); err != nil {
		return err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	t.iface = dev

	return nil
}

// Close closes the TAP device
func (t *TAPDevice) Close() error {
	return t.iface.Close()
}

// Read reads from the TAP device
func (t *TAPDevice) Read() error {
	for {
		frame := make([]byte, constants.FRAME_SIZE)

		_, err := t.iface.Read(frame)
		if err != nil {
			return err
		}

		t.framesChan <- frame
	}
}

// Write writes to the TAP device
func (t *TAPDevice) Write(frame []byte) (int, error) {
	return t.iface.Write(frame)
}
