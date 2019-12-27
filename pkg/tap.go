package pkg

import (
	"github.com/songgao/water"
	"os/exec"
)

type TAP struct {
	Name   string
	device water.Interface
}

func (t *TAP) Init(errors chan error, status chan string) {
	status <- "creating TAP device"

	config := water.Config{
		DeviceType: water.TAP,
	}
	config.Name = t.Name

	device, err := water.New(config)
	if err != nil {
		errors <- err
		return
	}

	t.device = *device

	status <- "created TAP device"

	status <- "bringing TAP device up"

	if _, err := exec.Command("ip", "link", "set", "dev", t.Name, "up").CombinedOutput(); err != nil {
		errors <- err
		return
	}

	status <- "brought TAP device up"

	// TODO: Close TAP device (in interrupt handler, maybe)
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

	status <- "read frames from TAP device"
}
