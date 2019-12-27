package device

import (
	"github.com/mdlayher/ethernet"
	"github.com/pojntfx/gloeth/pkg/protocol"
	ethernetRead "github.com/songgao/packets/ethernet"
	"github.com/songgao/water"
	"log"
	"net"
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

	var frame ethernetRead.Frame

	for {
		frame.Resize(1500)
		n, err := d.device.Read(frame)
		if err != nil {
			errors <- err
		}
		frame = frame[:n]

		source := net.HardwareAddr{frame.Source()[0], frame.Source()[1], frame.Source()[2], frame.Source()[3], frame.Source()[4], frame.Source()[5]}
		destination := net.HardwareAddr{frame.Destination()[0], frame.Destination()[1], frame.Destination()[2], frame.Destination()[3], frame.Destination()[4], frame.Destination()[5]}

		frameToSend := protocol.Frame{
			From: source.String(),
			To:   destination.String(),
			Body: frame.Payload(),
		}

		framesToSend <- frameToSend
	}
}
