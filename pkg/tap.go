package pkg

import (
	"github.com/songgao/water"
	"net"
	"os/exec"
)

type TAP struct {
	Name   string
	device water.Interface
}

func (t *TAP) Init() error {
	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = t.Name

	device, err := water.New(config)
	if err != nil {
		return err
	}

	t.device = *device

	if _, err := exec.Command("ip", "link", "set", "dev", t.Name, "up").CombinedOutput(); err != nil {
		return err
	}

	// TODO: Close TAP device (in interrupt handler, maybe)
	return nil
}

func (t *TAP) GetMacAddress() (error, string) {
	device, err := net.InterfaceByName(t.Name)
	if err != nil {
		return err, ""
	}

	return nil, device.HardwareAddr.String()
}

func (t *TAP) Write(errors chan error, status chan string, frame []byte) {
	status <- "writing frame to TAP device"

	if _, err := t.device.Write(frame); err != nil {
		errors <- err
		t.Write(errors, status, frame)
		return
	}

	status <- "wrote frame to TAP device"
}

func (t *TAP) Read(errors chan error, status chan string, readFrames chan []byte) {
	status <- "reading frames from TAP device"

	frame := make([]byte, 2000)

	for {
		status <- "reading frame from TAP device"

		n, err := t.device.Read(frame)
		if err != nil {
			errors <- err
		}

		readFrames <- frame[:n]

		status <- "read frame from TAP device"
	}
}
