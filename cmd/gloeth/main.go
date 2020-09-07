package main

import (
	"flag"
	"log"
	"sync"
)

func main() {
	// Parse flags
	deviceName := flag.String("deviceName", "gloeth0", "Network device name")
	maximumTransmissionUnit := flag.Int("maximumTransmissionUnit", 1500, "Frame size")

	preSharedKey := flag.String("preSharedKey", "supersecurekey", "Pre-shared key")

	localAddress := flag.String("localAddress", "0.0.0.0:1927", "Local address")
	localCertificate := flag.String("localCertificate", "/etc/gloeth/local.crt", "Local certificate")
	localKey := flag.String("localKey", "/etc/gloeth/local.key", "Local key")

	remoteAddress := flag.String("remoteAddress", "example.com:1927", "Remote address")
	remoteCertificate := flag.String("remoteCertificate", "/etc/gloeth/remote.crt", "Remote certificate")
	genesis := flag.Bool("genesis", false, "Enable genesis mode")

	debug := flag.Bool("debug", false, "Enable debugging mode")

	flag.Parse()

	// Create instances
	preSharedKeyValidator := validators.NewPreSharedKeyValidator(*preSharedKey)
	frameService := services.NewFrameService(preSharedKeyValidator)
	frameServer := servers.NewFrameServer(*localAddress, *localCertificate, *localKey, frameService)
	frameClient := clients.NewFrameClient(*remoteAddress, *remoteCertificate, preSharedKeyValidator)
	tapDevice := devices.NewTapDevice(*deviceName, *maximumTransmissionUnit)

	// Open instances
	if *debug {
		log.Println("opening instances")
	}

	if *genesis {
		if err := frameServer.Open(); err != nil {
			log.Fatal("could not open frame server", err)
		}
	} else {
		if err := frameClient.Open(); err != nil {
			log.Fatal("could not open frame client", err)
		}
	}

	if err := tapDevice.Open(); err != nil {
		log.Fatal("could not open TAP device", err)
	}

	// Connect instances
	var wg sync.WaitGroup

	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		for {
			frame, err := tapDevice.Read()
			if err != nil {
				log.Println("could not read from TAP device, dropping frame", err)
			}

			if *genesis {
				if err := frameService.Write(frame); err != nil {
					log.Println("could not write to frame service, dropping frame", err)
				}
			} else {
				if err := frameClient.Write(frame); err != nil {
					log.Println("could not write to frame client, dropping frame", err)
				}
			}
		}

		wg.Done()
	}(&wg)

	if *genesis {
		go func(wg *sync.WaitGroup) {
			for {
				frame, err := frameService.Read()
				if err != nil {
					log.Println("could not read from frame service, dropping frame", err)
				}

				if err := tapDevice.Write(frame); err != nil {
					log.Println("could not write to TAP device, dropping frame", err)
				}
			}

			wg.Done()
		}(&wg)
	} else {
		go func(wg *sync.WaitGroup) {
			for {
				frame, err := frameClient.Read()
				if err != nil {
					log.Println("could not read from frame client, dropping frame", err)
				}

				if err := tapDevice.Write(frame); err != nil {
					log.Println("could not write to TAP device, dropping frame", err)
				}
			}

			wg.Done()
		}(&wg)
	}

	wg.Wait()
}
