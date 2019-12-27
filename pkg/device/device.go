package device

import "github.com/pojntfx/gloeth/pkg/protocol"

type Device interface {
	Init() error
	Write(frame protocol.Frame) error
	Read(chan error, chan protocol.Frame)
}
