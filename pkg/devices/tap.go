package devices

import (
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

type TAPDevice struct {
	deviceName              string
	maximumTransmissionUnit int
	device                  *water.Interface
}

func NewTAPDevice(deviceName string, maximumTransmissionUnit int) *TAPDevice {
	return &TAPDevice{deviceName, maximumTransmissionUnit, nil}
}

func (d *TAPDevice) Open() error {
	device, err := water.New(water.Config{
		DeviceType: water.TAP,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: d.deviceName,
		},
	})
	if err != nil {
		return err
	}

	link, err := netlink.LinkByName(d.deviceName)
	if err != nil {
		return err
	}

	if err := netlink.LinkSetMTU(link, d.maximumTransmissionUnit); err != nil {
		return err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	d.device = device

	return nil
}

func (d *TAPDevice) Write(rawFrame []byte) error {
	_, err := d.device.Write(rawFrame)

	return err
}

func (d *TAPDevice) Read() ([]byte, error) {
	readFrame := make([]byte, d.maximumTransmissionUnit)

	_, err := d.device.Read(readFrame)

	return readFrame, err
}
