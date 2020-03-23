package main

import (
	"flag"
	"log"
	"net"

	"github.com/pojntfx/gloeth/pkg/connections"
	"github.com/pojntfx/gloeth/pkg/constants"
	"github.com/pojntfx/gloeth/pkg/devices"
	"github.com/pojntfx/gloeth/pkg/protocol"
)

func main() {
	localAddrFlag := flag.String("localAddr", "0.0.0.0:1234", "Local address")
	remoteAddrFlag := flag.String("remoteAddr", "0.0.0.0:12345", "Remote address")
	name := flag.String("name", "tap0", "Device name")
	flag.Parse()

	localAddr, err := net.ResolveTCPAddr("tcp", *localAddrFlag)
	if err != nil {
		log.Fatal(err)
	}

	remoteAddr, err := net.ResolveTCPAddr("tcp", *remoteAddrFlag)
	if err != nil {
		log.Fatal(err)
	}

	framesFromDeviceChan, framesFromConnectionChan := make(chan []byte), make(chan []byte)

	dev := devices.NewTAPDevice(constants.MTU, *name, framesFromDeviceChan)
	conn := connections.NewTAPviaTCPConnection(localAddr, remoteAddr, framesFromConnectionChan)
	enc := protocol.NewEncoder()

	defer dev.Close()
	if err := dev.Open(); err != nil {
		log.Fatal(err)
	}
	go dev.Read()

	defer conn.Close()
	if err := conn.Open(); err != nil {
		log.Fatal(err)
	}
	go conn.Read()

	go func() {
		for {
			inFrame := <-framesFromDeviceChan

			go func() {
				outFrame, invalid := enc.Encapsulate(inFrame)
				if invalid != nil {
					return
				}

				if _, err := conn.Write(outFrame); err != nil {
					log.Fatal(err)
				}
			}()
		}
	}()

	for {
		inFrame := <-framesFromConnectionChan

		go func() {
			outFrame, invalid := enc.Decapsulate(inFrame)
			if invalid != nil {
				return
			}

			if _, err := dev.Write(outFrame); err != nil {
				log.Fatal(err)
			}
		}()
	}
}
