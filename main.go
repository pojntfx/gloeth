package main

import (
	"flag"
	"github.com/pojntfx/gloeth/pkg"
	"log"
)

func main() {
	var (
		tapName          = flag.String("device", "gloeth", "Name of the TAP device to create")
		tcpReadHostPort  = flag.String("listen", "127.0.0.1:1234", "Host:port to listen on")
		tcpWriteHostPort = flag.String("peer", "127.0.0.1:1235", "Host:port the peer listens on")
	)
	flag.Parse()

	var (
		tapErrorChan = make(chan error)
		tcpErrorChan = make(chan error)

		tapStatusChan = make(chan string)
		tcpStatusChan = make(chan string)

		tapReadFramesChan = make(chan []byte)
		tcpReadFramesChan = make(chan []byte)
	)

	tap := pkg.TAP{
		Name: *tapName,
	}

	go tap.Init(tapErrorChan, tapStatusChan)

	for {
		var shouldBreak bool

		select {
		case err := <-tapErrorChan:
			log.Fatalln("TAP init error:", err)
		case status := <-tapStatusChan:
			log.Println("TAP init status:", status)
			if status == "brought TAP device up" {
				shouldBreak = true
			}
		}
		if shouldBreak {
			break
		}
	}

	tcp := pkg.TCP{
		WriteHostPort: *tcpWriteHostPort,
		ReadHostPort:  *tcpReadHostPort,
	}

	go tap.Read(tapErrorChan, tapStatusChan, tapReadFramesChan)
	go tcp.Read(tcpErrorChan, tcpStatusChan, tcpReadFramesChan)

	for {
		select {
		case err := <-tapErrorChan:
			log.Println("TAP error:", err)
		case err := <-tcpErrorChan:
			log.Println("TCP error:", err)

		case status := <-tapStatusChan:
			log.Println("TAP status:", status)
		case status := <-tcpStatusChan:
			log.Println("TCP status:", status)

		case frame := <-tapReadFramesChan:
			go tcp.Write(tcpErrorChan, tcpStatusChan, frame)
		case frame := <-tcpReadFramesChan:
			go tap.Write(tapErrorChan, tapStatusChan, frame)
		}
	}
}
