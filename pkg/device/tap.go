package device

import (
	"github.com/mdlayher/ethernet"
	"github.com/pojntfx/gloeth/pkg/protocol"
	ethernetRead "github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
	"log"
	"os/exec"
)

type TAP struct {
	Name   string
	device water.Interface
}

func (d *TAP) Init() error {
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = d.Name

	device, err := water.New(config)
	if err != nil {
		return err
	}

	d.device = *device

	if _, err := exec.Command("ip", "link", "set", "dev", d.Name, "up").CombinedOutput(); err != nil {
		return err
	}

	return nil
}

func (d *TAP) Write(frame protocol.Frame) error {
	etherFrame := &ethernet.Frame{
		Destination: []byte(frame.To),
		Source:      []byte(frame.From),
		EtherType:   0xcccc,
		Payload:     frame.Body,
	}

	etherFrameBinary, err := etherFrame.MarshalBinary()
	if err != nil {
		return err
	}

	if _, err := d.device.Write(etherFrameBinary); err != nil {
		return err
	}

	return nil
}

func (d *TAP) Read(errors chan error, framesToSend chan protocol.Frame) {
	log.Println("tap device reading")

	var ethernetFrame ethernetRead.Frame

	for {
		ethernetFrame.Resize(1500)

		n, err := d.device.Read(ethernetFrame)
		if err != nil {
			errors <- err
		}

		ethernetFrame = ethernetFrame[:n]

		frame := protocol.Frame{
			From: string(ethernetFrame.Source()),
			To:   string(ethernetFrame.Destination()),
			Body: ethernetFrame.Payload(),
		}

		framesToSend <- frame
	}
}
