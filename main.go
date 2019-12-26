package main

import (
	"flag"
	"github.com/pojntfx/gloeth/pkg/device"
	"github.com/pojntfx/gloeth/pkg/protocol"
	"github.com/pojntfx/gloeth/pkg/transceiver"
	"log"
)

func main() {
	tcpSendHostPort := flag.String("peer", "127.0.0.1:1235", "Host:port of the peer to send to")
	tcpListenHostPort := flag.String("listen", "127.0.0.1:1234", "Host:port to listen on")

	tapName := flag.String("device", "goeth", "Ethernet device to create")
	flag.Parse()

	errorsWhileReceivingFrames := make(chan error)
	errorsWhileSendingFrames := make(chan error)

	receivedFrames := make(chan protocol.Frame)
	framesToSend := make(chan protocol.Frame)

	tcp := transceiver.TCP{
		SendHostPort:   *tcpSendHostPort,
		ListenHostPort: *tcpListenHostPort,
	}
	go tcp.Listen(errorsWhileReceivingFrames, receivedFrames)

	tap := device.TAP{
		Name: *tapName,
	}
	go tap.Read(errorsWhileSendingFrames, framesToSend)

	for {
		select {
		case err := <-errorsWhileReceivingFrames:
			log.Println(err)
		case err := <-errorsWhileSendingFrames:
			log.Println(err)

		case frame := <-receivedFrames:
			if err := tap.Write(frame); err != nil {
				errorsWhileReceivingFrames <- err
			}
		case frame := <-framesToSend:
			if err := tcp.Send(frame); err != nil {
				errorsWhileSendingFrames <- err
			}
		}
	}
}
