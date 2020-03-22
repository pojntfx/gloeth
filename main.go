// Based on https://github.com/vsergeev/tinytaptunnel

package main

import (
	"flag"
	"log"
	"net"

	"github.com/pojntfx/gloeth/pkg/utils"
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

	tcpListener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		log.Fatalf("error creating a TCP socket: %v\n", err)
	}

	tapDev := new(utils.TAPDevice)
	err = tapDev.Open(utils.TAP_MTU)
	if err != nil {
		log.Fatalf("could not open a TAP device: %v\n", err)
	}

	log.Printf("started tunnel with TAP device %v", tapDev.GetName())

	go utils.ForwardTCPtoTAP(tcpListener, tapDev, remoteAddr)
	utils.ForwardTAPtoTCP(tapDev, remoteAddr)
}
