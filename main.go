package main

import (
	"flag"
	"log"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:1234", "Host:port to listen on")
	peer := flag.String("peer", "127.0.0.1:1235", "Host:port to connect to")
	device := flag.String("device", "goeth", "Ethernet device to create")
	flag.Parse()

	log.Println(*listen, "=>", *device, "=>", *peer)
}
