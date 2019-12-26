package device

import (
	"github.com/pojntfx/gloeth/pkg/protocol"
	"log"
)

type TAP struct {
	Name string
}

func (d TAP) Write(frame protocol.Frame) error {
	log.Println("tap device writing frame", frame)

	return nil
}

func (d TAP) Read(errors chan error, framesToSend chan protocol.Frame) {
	log.Println("tap device reading")
}
