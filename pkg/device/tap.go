package device

import (
	"github.com/pojntfx/gloeth/pkg/protocol"
	"github.com/songgao/water"
	"log"
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

	return nil
}

func (d TAP) Write(frame protocol.Frame) error {
	log.Println("tap device writing frame", frame)

	return nil
}

func (d TAP) Read(errors chan error, framesToSend chan protocol.Frame) {
	log.Println("tap device reading")
}
