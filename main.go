package main

import (
	"flag"
	"github.com/pojntfx/gloeth/pkg/device"
	"github.com/pojntfx/gloeth/pkg/protocol"
	"github.com/pojntfx/gloeth/pkg/transceiver"
)

func main() {
	tcpSendHostPort := flag.String("peer", "127.0.0.1:1235", "Host:port of the peer to send to")
	tcpListenHostPort := flag.String("listen", "127.0.0.1:1234", "Host:port to listen on")

	tapName := flag.String("device", "goeth", "Ethernet device to create")
	flag.Parse()

	receivedFrames := make(chan protocol.Frame)
	framesToSend := make(chan protocol.Frame)

	tcp := transceiver.TCP{
		SendHostPort:   *tcpSendHostPort,
		ListenHostPort: *tcpListenHostPort,
	}
	tcp.Listen(&receivedFrames)

	tap := device.TAP{
		Name: *tapName,
	}
	tap.Listen(&framesToSend)

	for {
		select {
		case frame := <-receivedFrames:
			tap.Write(frame)
		case frame := <-framesToSend:
			tcp.Send(frame)
		}
	}
}
