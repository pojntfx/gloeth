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

	initConns := make(map[string]*net.TCPConn)
	switcher := switchers.NewTCP(readChan, connChan, laddr, initConns)

	defer switcher.Close()
	if err := switcher.Open(); err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %v", laddr)

	wpr := wrappers.NewEthernet()

	go func() {
		if err := switcher.Read(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		for {
			conn := <-connChan

			go func() {
				for {
					inFrame := [wrappers.WrappedFrameSize]byte{}

					_, err := conn.Read(inFrame[:])
					if err != nil {
						log.Printf("could not read from connection: %v", err)

						break
					}

					_, sourceMAC, _, err := wpr.Unwrap(inFrame)
					if err != nil {
						log.Printf("could not unwrap frame: %v", err)

						continue
					}

					switcher.Register(sourceMAC, conn)

					switcher.HandleFrame(inFrame)
				}
			}()
		}
	}()

	for {
		inFrame := <-readChan

		destMAC, sourceMAC, _, err := wpr.Unwrap(inFrame)
		if err != nil {
			log.Printf("could not unwrap frame: %v", err)

			continue
		}

		conns, err := switcher.GetConnectionsForMAC(destMAC, sourceMAC)
		if err != nil {
			log.Printf("could not get connections: %v", err)

			continue
		}

		for _, conn := range conns {
			if err := switcher.Write(conn, inFrame); err != nil {
				log.Printf("could not write to connection: %v", err)
			}
		}
	}
}
