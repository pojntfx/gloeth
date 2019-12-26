package transceiver

import (
	"github.com/pojntfx/gloeth/pkg/protocol"
)

type TCP struct {
	SendHostPort   string
	ListenHostPort string
}

func (t TCP) Send(frame protocol.Frame) error {
	panic("implement me")
}

func (t TCP) Listen(*chan protocol.Frame) error {
	panic("implement me")
}
