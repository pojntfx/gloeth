package main

import (
	"flag"
	"log"
	"net"

	"github.com/pojntfx/gloeth/v3/pkg/connections"
	"github.com/pojntfx/gloeth/v3/pkg/devices"
	"github.com/pojntfx/gloeth/v3/pkg/encoders"
	"github.com/pojntfx/gloeth/v3/pkg/encryptors"
	"github.com/pojntfx/gloeth/v3/pkg/wrappers"
)

func main() {
	saddrFlag := flag.String("saddr", ":1234", "Supernode address")
	key := flag.String("key", "my_preshared_key", "Preshared key")
	name := flag.String("name", "tap0", "Device name")
	flag.Parse()

	saddr, err := net.ResolveTCPAddr("tcp", *saddrFlag)
	if err != nil {
		log.Fatal(err)
	}

	devChan, connChan := make(chan [encryptors.PlaintextFrameSize]byte), make(chan [wrappers.WrappedFrameSize]byte)

	dev := devices.NewTAP(devChan, devices.MTU, *name)
	conn := connections.NewTCPSwitcher(connChan, saddr)

	defer dev.Close()
	if err := dev.Open(); err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	if err := conn.Open(); err != nil {
		log.Fatal(err)
	}

	enco := encoders.NewEthernet()
	encr := encryptors.NewEthernet(*key)
	wpr := wrappers.NewEthernet()

	go func() {
		for {
			inFrame := <-devChan

			addr, err := enco.GetDestMACAddress(inFrame)
			if err != nil {
				log.Println(err)

				continue
			}

			encrFrame, err := encr.Encrypt(inFrame)
			if err != nil {
				log.Println(err)

				continue
			}

			outFrame, err := wpr.Wrap(addr, encrFrame)
			if err != nil {
				log.Println(err)

				continue
			}

			if err := conn.Write(outFrame); err != nil {
				log.Println(err)

				continue
			}
		}
	}()

	for {
		inFrame := <-connChan

		_, dewrpFrame, err := wpr.Unwrap(inFrame)
		if err != nil {
			log.Println(err)

			continue
		}

		decrFrame, err := encr.Decrypt(dewrpFrame)
		if err != nil {
			log.Println(err)

			continue
		}

		if err := dev.Write(decrFrame); err != nil {
			log.Println(err)

			continue
		}
	}
}
