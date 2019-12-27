package device

import (
	"github.com/pojntfx/gloeth/pkg/protocol"
	"github.com/songgao/water"
	"log"
	"os/exec"
)

type TAP struct {
	Name   string
	device *water.Interface
}

func (d TAP) Init() error {
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = d.Name

	device, err := water.New(config)
	if err != nil {
		return err
	}

	d.device = device

	if _, err := exec.Command("ip", "link", "set", "dev", d.Name, "up").CombinedOutput(); err != nil {
		return err
	}

	return nil
}

func (d TAP) Write(frame protocol.Frame) error {
	return nil
}

func (d TAP) Read(errors chan error, framesToSend chan protocol.Frame) {
	log.Println("tap device reading")
}
