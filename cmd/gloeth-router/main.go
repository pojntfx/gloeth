package main

import (
	"bytes"
	"flag"
	"log"
	"net"

	"github.com/pojntfx/gloeth/pkg/constants"
	"github.com/pojntfx/gloeth/pkg/protocol"
	"github.com/songgao/packets/ethernet"
)

var conns map[string]*net.TCPConn

func main() {
	listenAddrFlag := flag.String("listenAddr", ":1234", "The listen address")
	flag.Parse()

	listenAddr, err := net.ResolveTCPAddr("tcp", *listenAddrFlag)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}

	enc := protocol.NewEncoder()

	conn, err := l.AcceptTCP()
	if err != nil {
		log.Fatal(err)
	}

	for {
		frame := make([]byte, constants.FRAME_SIZE+constants.TIMESTAMP_SIZE)

		_, err = conn.Read(frame)
		if err != nil {
			log.Fatal(err)
		}

		decFrame, err := enc.Decapsulate(frame)
		if err != nil {
			continue
		}

		var ethernetFrame ethernet.Frame
		ethernetFrame.Resize(constants.FRAME_SIZE)

		r := bytes.NewReader(decFrame)
		if _, err := r.Read(ethernetFrame); err != nil {
			log.Fatal(err)
		}

		log.Println(ethernetFrame.Source(), ethernetFrame.Destination(), ethernetFrame.Ethertype(), string(ethernetFrame.Payload()))

		// Now write to the one other connection
		// if _, err := conn.Write(frame); err != nil {
		// log.Fatal(err)
		// }
	}
}
