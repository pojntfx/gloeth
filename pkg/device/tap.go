package device

import (
	"github.com/pojntfx/gloeth/pkg/protocol"
)

type TAP struct {
	Name string
}

func (d TAP) Write(frame protocol.Frame) error {
	panic("implement me")
}

func (d TAP) Listen(*chan protocol.Frame) error {
	panic("implement me")
}
