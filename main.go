// Based on https://github.com/vsergeev/tinytaptunnel

package main

import (
	"flag"
	"log"
	"net"

	"github.com/pojntfx/gloeth/pkg/constants"
	"github.com/pojntfx/gloeth/pkg/forwarding"
	"github.com/pojntfx/gloeth/pkg/switcher"
	"github.com/pojntfx/gloeth/pkg/tap"
)

func main() {
	localAddrFlag := flag.String("localAddr", "0.0.0.0:1234", "Local address")
	remoteAddrFlag := flag.String("remoteAddr", "0.0.0.0:12345", "Remote address")

	flag.Parse()

	localAddr, err := net.ResolveTCPAddr("tcp", *localAddrFlag)
	if err != nil {
		log.Fatalf("could not resolve local address: %v\n", err)
	}

	remoteAddr, err := net.ResolveTCPAddr("tcp", *remoteAddrFlag)
	if err != nil {
		log.Fatalf("could not resolve remote address: %v\n", err)
	}

	tapDev := tap.NewDevice()
	err = tapDev.Open(constants.TAP_MTU)
	if err != nil {
		log.Fatalf("could not open a TAP device: %v\n", err)
	}

	log.Printf("started tunnel with TAP device %v", tapDev.GetName())

	switcherConnection := switcher.NewConnection(localAddr, remoteAddr)
	errChan := make(chan error)

	forwarder := forwarding.NewTAPviaTCPForwarder(tapDev, localAddr, remoteAddr, switcherConnection, errChan)

	go func() {
		for err := range errChan {
			log.Println(err)
		}
	}()

	go forwarder.TCPtoTAP()
	forwarder.TAPtoTCP()
}
