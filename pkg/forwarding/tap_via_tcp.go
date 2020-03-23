package forwarding

import (
	"fmt"
	"net"

	"github.com/pojntfx/gloeth/pkg/constants"
	"github.com/pojntfx/gloeth/pkg/encoding"
	"github.com/pojntfx/gloeth/pkg/switcher"
	"github.com/pojntfx/gloeth/pkg/tap"
)

// TAPviaTCP forwards TCP to TAP and vice versa
type TAPviaTCP struct {
	tapDevice             *tap.Device
	localAddr, remoteAddr *net.TCPAddr
	switcherConnection    *switcher.Connection
	errChan               chan error
}

// NewTAPviaTCPForwarder creates a new TAP via TCP forwarder
func NewTAPviaTCPForwarder(tapDevice *tap.Device, localAddr, remoteAddr *net.TCPAddr, switcherConnection *switcher.Connection, errChan chan error) *TAPviaTCP {
	return &TAPviaTCP{tapDevice: tapDevice, localAddr: localAddr, remoteAddr: remoteAddr, switcherConnection: switcherConnection, errChan: errChan}
}

// TCPtoTAP forwards TCP packets to a TAP device
func (f *TAPviaTCP) TCPtoTAP() {
	readResChan := make(chan switcher.ReadPacket)
	go func() {
		if err := f.switcherConnection.Read(readResChan); err != nil {
			f.errChan <- fmt.Errorf("could not read from switcher connection %v", err)
		}
	}()

	for {
		readRes := <-readResChan

		frame, invalid := encoding.DecapsulateFrame(readRes.Payload[0:readRes.Length])
		if invalid != nil {
			continue
		}

		_, err := f.tapDevice.Write(frame)
		if err != nil {
			f.errChan <- fmt.Errorf("could not write to TAP device: %v", err)

			return
		}
	}
}

// TAPtoTCP forwards frames from a TAP device to a TCP connection
func (f *TAPviaTCP) TAPtoTCP() {
	for {
		frame := make([]byte, constants.TAP_MTU+14)
		var encFrame []byte

		n, err := f.tapDevice.Read(frame)
		if err != nil {
			f.errChan <- fmt.Errorf("could not read from TAP device: %v", err)
		}

		encFrame, invalid := encoding.EncapsulateFrame(frame[0:n])
		if invalid != nil {
			continue
		}

		if err := f.switcherConnection.Write(encFrame); err != nil {
			f.errChan <- fmt.Errorf("could not dial, retrying: %v", err)

			f.TAPtoTCP()

			return
		}
	}
}
