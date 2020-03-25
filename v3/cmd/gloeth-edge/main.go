package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/pojntfx/gloeth/v3/pkg/connections"
	"github.com/pojntfx/gloeth/v3/pkg/devices"
	"github.com/pojntfx/gloeth/v3/pkg/encoders"
	"github.com/pojntfx/gloeth/v3/pkg/encryptors"
	"github.com/pojntfx/gloeth/v3/pkg/wrappers"
)

func main() {
	raddrFlag := flag.String("raddr", ":1234", "Supernode address")
	key := flag.String("key", "my_preshared_key", "Preshared key")
	name := flag.String("name", "tap0", "Device name")
	flag.Parse()

	raddr, err := net.ResolveTCPAddr("tcp", *raddrFlag)
	if err != nil {
		log.Fatal(err)
	}

	devChan, connChan := make(chan [encryptors.PlaintextFrameSize]byte), make(chan [wrappers.WrappedFrameSize]byte)

	dev := devices.NewTAP(devChan, devices.MTU, *name)
	conn := connections.NewTCPSwitcher(connChan, raddr)

	defer dev.Close()
	if err := dev.Open(); err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	if err := conn.Open(); err != nil {
		log.Fatal(err)
	}
	log.Printf("successfully connected to supernode %v", raddr)

	enco := encoders.NewEthernet()
	encr := encryptors.NewEthernet(*key)
	wpr := wrappers.NewEthernet()

	go func() {
		if err := dev.Read(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		timeTillReconnect := time.Millisecond * 250

		for {
			if err := conn.Read(); err != nil {
				log.Printf("could not read from supernode %v due to error %v, retrying now", raddr, err)
			}

			if err := conn.Open(); err != nil {
				log.Printf("could not reconnect from supernode %v due to error %v, retrying in %v", raddr, err, timeTillReconnect)

				time.Sleep(timeTillReconnect)

				continue
			}

			log.Printf("successfully reconnected to supernode %v", raddr)
		}
	}()

	go func() {
		for {
			inFrame := <-devChan

			destMAC, srcMAC, err := enco.GetMACAddresses(inFrame)
			if err != nil {
				log.Printf("could not get MAC addresses from ethernet frame: %v", err)

				continue
			}

			encrFrame, err := encr.Encrypt(inFrame)
			if err != nil {
				log.Printf("could not encrypt ethernet frame: %v", err)

				continue
			}

			outFrame, err := wpr.Wrap(destMAC, srcMAC, encrFrame)
			if err != nil {
				log.Printf("could not wrap frame: %v", err)

				continue
			}

			if err := conn.Write(outFrame); err != nil {
				log.Printf("could not write frame to supernode: %v", err)

				continue
			}
		}
	}()

	for {
		inFrame := <-connChan

		_, _, dewrpFrame, err := wpr.Unwrap(inFrame)
		if err != nil {
			log.Printf("could not unwrap frame: %v", err)

			continue
		}

		decrFrame, err := encr.Decrypt(dewrpFrame)
		if err != nil {
			log.Printf("could not decrypt ethernet frame: %v", err)

			continue
		}

		if err := dev.Write(decrFrame); err != nil {
			log.Printf("could not write ethernet frame to device: %v", err)

			continue
		}
	}
}
