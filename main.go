package main

import (
	"flag"
	"github.com/pojntfx/gloeth/pkg/device"
	"github.com/pojntfx/gloeth/pkg/protocol"
	"github.com/pojntfx/gloeth/pkg/transceiver"
	"log"
)

func main() {
	tapName := flag.String("device", "goeth", "Ethernet device to create")

	tcpSendHostPort := flag.String("peer", "127.0.0.1:1235", "Host:port of the peer to send to")
	tcpListenHostPort := flag.String("listen", "127.0.0.1:1234", "Host:port to listen on")
	flag.Parse()

	errorsWhileSendingFrames := make(chan error)
	errorsWhileReceivingFrames := make(chan error)

	framesToSend := make(chan protocol.Frame)
	receivedFrames := make(chan protocol.Frame)

	tap := device.TAP{
		Name: *tapName,
	}
	if err := tap.Init(); err != nil {
		log.Fatal(err)
	}
	go tap.Read(errorsWhileSendingFrames, framesToSend)

	tcp := transceiver.TCP{
		SendHostPort:   *tcpSendHostPort,
		ListenHostPort: *tcpListenHostPort,
	}
	go tcp.Listen(errorsWhileReceivingFrames, receivedFrames)

	for {
		select {
		case err := <-errorsWhileSendingFrames:
			log.Println(err)
		case err := <-errorsWhileReceivingFrames:
			log.Println(err)

		case frame := <-framesToSend:
			log.Println("Read frame from TAP device, sending with TCP transceiver", frame.From, frame.To, string(frame.Body))
			go tcp.Send(frame) // TODO: Handle error
			//if err := tcp.Send(frame); err != nil {
			//	errorsWhileSendingFrames <- err
			//}
		case frame := <-receivedFrames:
			log.Println("Received frame from TCP transceiver, writing to TAP device", frame.From, frame.To, string(frame.Body))
			go tap.Write(frame) // TODO: Handle error
			//if err := tap.Write(frame); err != nil {
			//	errorsWhileReceivingFrames <- err
			//}
		}
	}
}
