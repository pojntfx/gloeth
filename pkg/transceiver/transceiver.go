package transceiver

import "github.com/pojntfx/gloeth/pkg/protocol"

type Transceiver interface {
	Send(frame protocol.Frame) error
	Listen(*chan protocol.Frame) error
}
