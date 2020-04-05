package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/pojntfx/gloeth/pkg/connections"
	"github.com/pojntfx/gloeth/pkg/devices"
	"github.com/pojntfx/gloeth/pkg/encoders"
	"github.com/pojntfx/gloeth/pkg/encryptors"
	"github.com/pojntfx/gloeth/pkg/wrappers"
)

func main() {
	raddrFlag := flag.String("raddr", ":1234", "Remote address")
	key := flag.String("key", "my_preshared_key", "Preshared key")
	name := flag.String("name", "tap0", "Device name")
	verbose := flag.Bool("verbose", false, "Enable verbose mode")
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
	log.Printf("successfully connected to switcher %v", raddr)

	enco := encoders.NewEthernet()
	encr := encryptors.NewEthernet(*key)
	wpr := wrappers.NewEthernet()

	go func() {
		for {
			if err := dev.Read(); err != nil {
				log.Printf("could not read from dev: %v", err)
			}
		}
	}()

	go func() {
		timeTillReconnect := time.Millisecond * 250

		for {
			if err := conn.Read(); err != nil {
				log.Printf("could not read from switcher %v due to error %v, retrying now", raddr, err)
			}

			if err := conn.Open(); err != nil {
				log.Printf("could not reconnect to switcher %v due to error %v, retrying in %v", raddr, err, timeTillReconnect)

				time.Sleep(timeTillReconnect)

				continue
			}

			log.Printf("successfully reconnected to switcher %v", raddr)
		}
	}()

	go func() {
		for {
			inFrame := <-devChan

			if *verbose {
				log.Printf("READ frame from TAP device: %v", inFrame)
			}

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

			outFrame, err := wpr.Wrap(destMAC, srcMAC, [wrappers.HopsCount]*net.HardwareAddr{}, encrFrame)
			if err != nil {
				log.Printf("could not wrap frame: %v", err)

				continue
			}

			if *verbose {
				log.Printf("WRITING frame to switcher: %v", outFrame)
			}

			if err := conn.Write(outFrame); err != nil {
				log.Printf("could not write frame to switcher: %v", err)

				continue
			}
		}
	}()

	for {
		inFrame := <-connChan

		if *verbose {
			log.Printf("READ frame from switcher: %v", inFrame)
		}

		_, _, _, dewrpFrame, err := wpr.Unwrap(inFrame)
		if err != nil {
			log.Printf("could not unwrap frame: %v", err)

			continue
		}

		decrFrame, err := encr.Decrypt(dewrpFrame)
		if err != nil {
			log.Printf("could not decrypt ethernet frame: %v", err)

			continue
		}

		if *verbose {
			log.Printf("WRITING frame to TAP device: %v", decrFrame)
		}

		if err := dev.Write(decrFrame); err != nil {
			log.Printf("could not write ethernet frame to device: %v", err)

			continue
		}
	}
}
