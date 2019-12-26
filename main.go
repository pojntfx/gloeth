package main

import (
	"flag"
	server2 "github.com/pojntfx/gloeth/pkg/server"
	"log"
	"net"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:1234", "Host:port to listen on")
	peer := flag.String("peer", "127.0.0.1:1235", "Host:port to connect to")
	device := flag.String("device", "goeth", "Ethernet device to create")
	flag.Parse()

	server, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalln("Could not listen", err)
	}
	defer server.Close()

	log.Println(*listen, "=>", *device, "=>", *peer)

	for {
		connection, err := server.Accept()
		if err != nil {
			log.Fatalln("Could not accept connection", err)
		}

		go server2.HandleConnection(connection)
	}
}
