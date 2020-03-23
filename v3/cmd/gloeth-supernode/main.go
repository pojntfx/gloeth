package main

import (
	"flag"
	"log"
	"net"

	"github.com/pojntfx/gloeth/v3/pkg/switchers"
	"github.com/pojntfx/gloeth/v3/pkg/wrappers"
)

func main() {
	laddrFlag := flag.String("laddr", ":1234", "Listen address")
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *laddrFlag)
	if err != nil {
		log.Fatal(err)
	}

	readChan := make(chan []byte)

	switcher := switchers.NewTCP(readChan, laddr)

	defer switcher.Close()
	if err := switcher.Open(); err != nil {
		log.Fatal(err)
	}

	wpr := wrappers.NewEthernet()

	for {
		inFrame := <-readChan

		unwrpFrame, addr, err := wpr.Unwrap(inFrame)
		if err != nil {
			log.Println(err)

			continue
		}

		conns, err := switcher.GetConnectionsForMAC(addr)
		if err != nil {
			log.Println(err)

			continue
		}

		for _, conn := range conns {
			if _, err := conn.Write(unwrpFrame); err != nil {
				log.Println(err)
			}
		}
	}
}
