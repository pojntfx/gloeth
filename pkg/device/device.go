package device

import "github.com/pojntfx/gloeth/pkg/protocol"

type Device interface {
	Write(frame protocol.Frame) error
	Listen(*chan protocol.Frame) error
}
