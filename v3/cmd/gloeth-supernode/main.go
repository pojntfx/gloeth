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

	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	connChan := make(chan *net.TCPConn)

	conns := make(map[string]*net.TCPConn)
	switcher := switchers.NewTCP(readChan, connChan, laddr, conns)

	defer switcher.Close()
	if err := switcher.Open(); err != nil {
		log.Fatal(err)
	}

	wpr := wrappers.NewEthernet()

	go func() {
		if err := switcher.Read(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		conn := <-connChan

		go func() {
			for {
				inFrame := [wrappers.WrappedFrameSize]byte{}

				_, err := conn.Read(inFrame[:])
				if err != nil {
					log.Fatal(err)
				}

				destMAC, _, _, err := wpr.Unwrap(inFrame)
				if err != nil {
					log.Println(err)

					continue
				}

				switcher.Register(destMAC, conn)

				switcher.HandleFrame(inFrame)
			}
		}()
	}()

	for {
		inFrame := <-readChan

		destMAC, sourceMAC, _, err := wpr.Unwrap(inFrame)
		if err != nil {
			log.Println(err)

			continue
		}

		conns, err := switcher.GetConnectionsForMAC(destMAC, sourceMAC)
		if err != nil {
			log.Println(err)

			continue
		}

		for _, conn := range conns {
			if err := switcher.Write(conn, inFrame); err != nil {
				log.Println(err)
			}
		}
	}
}
