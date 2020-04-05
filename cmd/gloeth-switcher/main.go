package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/pojntfx/gloeth/pkg/connections"
	"github.com/pojntfx/gloeth/pkg/switchers"
	"github.com/pojntfx/gloeth/pkg/wrappers"
)

func main() {
	laddrFlag := flag.String("laddr", ":1234", "Listen address")
	liaddrFlag := flag.String("liaddr", ":1235", "Listen address for info endpoint")
	raddrFlag := flag.String("raddr", ":1236", "Remote address")
	riaddrFlag := flag.String("riaddr", ":1237", "Remote address for info endpoint")
	verbose := flag.Int("verbose", 0, "Enable verbose mode (1 = excluding frames, 2 = including frames)")
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *laddrFlag)
	if err != nil {
		log.Fatal(err)
	}
	liaddr, err := net.ResolveTCPAddr("tcp", *liaddrFlag)
	if err != nil {
		log.Fatal(err)
	}
	raddr, err := net.ResolveTCPAddr("tcp", *raddrFlag)
	if err != nil {
		log.Fatal(err)
	}
	riaddr, err := net.ResolveTCPAddr("tcp", *riaddrFlag)
	if err != nil {
		log.Fatal(err)
	}

	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	connChan := make(chan *net.TCPConn)
	remoteSwitcherInfoChan := make(chan *net.HardwareAddr)

	initConns := make(map[string]*net.TCPConn)
	switcher := switchers.NewTCP(readChan, connChan, laddr, initConns)
	switcherInfo := switchers.NewSwitcherInfo(liaddr)
	remoteSwitcherInfoConn := connections.NewSwitcherInfo(remoteSwitcherInfoChan, riaddr)

	defer switcher.Close()
	if err := switcher.Open(); err != nil {
		log.Fatal(err)
	}
	log.Printf("listening on %v", laddr)

	defer switcherInfo.Close()
	if err := switcherInfo.RequestMACAddress(); err != nil {
		log.Fatal(err)
	}
	if err := switcherInfo.Open(); err != nil {
		log.Fatal(err)
	}
	log.Printf("info endpoint listening on %v", liaddr)

	wpr := wrappers.NewEthernet()

	go func() {
		for {
			if err := switcher.Read(); err != nil {
				log.Printf("could not read from switcher: %v", err)
			}
		}
	}()

	go func() {
		if err := switcherInfo.Read(); err != nil {
			log.Printf("could not read from switcher info: %v", err)
		}
	}()

	go func() {
		timeTillReconnect := time.Millisecond * 250

		for {
			defer remoteSwitcherInfoConn.Close()
			if err := remoteSwitcherInfoConn.Open(); err != nil {
				log.Printf("could not connect to remote switcher info %v due to error %v, retrying in %v", riaddr, err, timeTillReconnect)

				time.Sleep(timeTillReconnect)

				continue
			}

			log.Printf("successfully connected to remote switcher info %v", raddr)

			if err := remoteSwitcherInfoConn.Read(); err != nil {
				log.Printf("could not read from remote switcher info %v due to error %v, retrying now", riaddr, err)
			}

			defer remoteSwitcherInfoConn.Close()
			if err := remoteSwitcherInfoConn.Open(); err != nil {
				log.Printf("could not reconnect to remote switcher info %v due to error %v, retrying in %v", riaddr, err, timeTillReconnect)

				time.Sleep(timeTillReconnect)

				continue
			}

			log.Printf("successfully reconnected to remote switcher info %v", raddr)
		}
	}()

	go func() {
		for {
			mac := <-remoteSwitcherInfoChan

			conn, err := connections.GetConn(raddr)
			if err != nil {
				log.Printf("could not connect to remote switcher %v with MAC address %v: %v", raddr, mac, err)
			}

			if *verbose > 0 {
				log.Printf("REGISTERING connection for switcher with MAC %v: %v", mac, conn)
			}

			switcher.Register(mac, conn)

			log.Printf("successfully connected to remote switcher %v with MAC address %v", raddr, mac)
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

					if *verbose > 1 {
						log.Printf("READ frame from adapter: %v", inFrame)
					}

					destMAC, sourceMAC, hops, _, err := wpr.Unwrap(inFrame)
					if err != nil {
						log.Printf("could not unwrap frame: %v", err)

						continue
					}

					if *verbose > 0 {
						log.Printf("READ hops for frame from %v to %v: %v", sourceMAC, destMAC, hops)
					}

					if *verbose > 0 {
						log.Printf("REGISTERING connection for adapter with MAC %v: %v", sourceMAC, conn)
					}

					switcher.Register(sourceMAC, conn)

					switcher.HandleFrame(inFrame)
				}
			}()
		}
	}()

	for {
		inFrame := <-readChan

		destMAC, sourceMAC, hops, frame, err := wpr.Unwrap(inFrame)
		if err != nil {
			log.Printf("could not unwrap frame: %v", err)

			continue
		}

		if *verbose > 0 {
			log.Printf("READ hops for frame from %v to %v: %v", sourceMAC, destMAC, hops)
		}

		newHops := wpr.GetShiftedHops(hops)
		if wpr.GetHopsEmpty(newHops) {
			conns, err := switcher.GetConnectionsForMAC(destMAC, sourceMAC)
			if err != nil {
				log.Printf("could not get connections: %v", err)

				continue
			}

			if *verbose > 0 {
				if *verbose > 1 {
					log.Printf("WRITING frame to adapter(s) with MAC %v via connections %v: %v", destMAC, conns, inFrame)
				} else {
					log.Printf("WRITING frame to adapter(s) with MAC %v via connections %v", destMAC, conns)
				}
			}

			for _, conn := range conns {
				if err := switcher.Write(conn, inFrame); err != nil {
					log.Printf("could not write to connection: %v", err)
				}
			}

			continue
		}

		newFrame, err := wpr.Wrap(destMAC, sourceMAC, newHops, frame)
		if err != nil {
			log.Printf("could not wrap frame: %v", err)

			continue
		}

		conns, err := switcher.GetConnectionsForMAC(hops[len(hops)-1], sourceMAC)
		if err != nil {
			log.Printf("could not get connections: %v", err)

			continue
		}

		if *verbose > 0 {
			if *verbose > 1 {
				log.Printf("WRITING frame to adapter(s) with MAC %v via connections %v: %v", destMAC, conns, inFrame)
			} else {
				log.Printf("WRITING frame to adapter(s) with MAC %v via connections %v", destMAC, conns)
			}
		}

		for _, conn := range conns {
			if err := switcher.Write(conn, newFrame); err != nil {
				log.Printf("could not write to connection: %v", err)
			}
		}

		continue
	}
}
