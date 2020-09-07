package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/pojntfx/gloeth/pkg/clients"
	"github.com/pojntfx/gloeth/pkg/converters"
	"github.com/pojntfx/gloeth/pkg/devices"
	"github.com/pojntfx/gloeth/pkg/servers"
	"github.com/pojntfx/gloeth/pkg/services"
	"github.com/pojntfx/gloeth/pkg/validators"
)

func main() {
	// Parse flags
	deviceName := flag.String("deviceName", "gloeth0", "Network device name")
	maximumTransmissionUnit := flag.Int("maximumTransmissionUnit", 1500, "Frame size")

	preSharedKey := flag.String("preSharedKey", "supersecurekey", "Pre-shared key")
	genesis := flag.Bool("genesis", false, "Enable genesis mode")

	localAddress := flag.String("localAddress", "0.0.0.0:1927", "Local address (only required when in genesis mode)")
	localCertificate := flag.String("localCertificate", "/etc/gloeth/local.crt", "Local certificate (only required when in genesis mode)")
	localKey := flag.String("localKey", "/etc/gloeth/local.key", "Local key (only required when in genesis mode)")

	remoteAddress := flag.String("remoteAddress", "example.com:1927", "Remote address (not required when in genesis mode)")
	remoteCertificate := flag.String("remoteCertificate", "/etc/gloeth/remote.crt", "Remote certificate (not required when in genesis mode)")

	debug := flag.Bool("debug", false, "Enable debugging mode")

	flag.Parse()

	// Create instances
	preSharedKeyValidator := validators.NewPreSharedKeyValidator(*preSharedKey)
	frameConverter := converters.NewFrameConverter()
	frameService := services.NewFrameService()
	frameServer := servers.NewFrameServer(*localAddress, *localCertificate, *localKey, frameService)
	frameClient := clients.NewFrameClient(*remoteAddress, *remoteCertificate)
	tapDevice := devices.NewTAPDevice(*deviceName, *maximumTransmissionUnit)

	// Open instances
	if *genesis {
		go func() {
			log.Println("Opening frame server")

			if err := frameServer.Open(); err != nil {
				log.Fatal("could not open frame server", err)
			}
		}()
	} else {
		go func() {
			log.Println("Opening frame client")

			if err := frameClient.Open(); err != nil {
				log.Fatal("could not open frame client", err)
			}
		}()
	}

	go func() {
		log.Println("Opening TAP device")

		if err := tapDevice.Open(); err != nil {
			log.Fatal("could not open TAP device", err)
		}
	}()

	// Connect instances
	var wg sync.WaitGroup

	wg.Add(2)

	go func(wg *sync.WaitGroup) {
		log.Println("Reading from TAP device")

		for {
			rawFrame, err := tapDevice.Read()
			if err != nil {
				log.Println("could not read from TAP device, dropping frame", err)

				continue
			}

			frame, err := frameConverter.ToExternal(rawFrame, *preSharedKey)
			if err != nil {
				log.Println("could not convert internal frame to external frame, dropping frame", err)

				continue
			}

			if *genesis {
				if *debug {
					log.Println("Writing frame from TAP device to frame service")
				}

				if err := frameService.Write(frame); err != nil {
					log.Println("could not write to frame service, dropping frame and continuing in 250ms", err)

					time.Sleep(time.Millisecond * 250)

					continue
				}
			} else {
				if *debug {
					log.Println("Writing frame from TAP device to frame client")
				}

				if err := frameClient.Write(frame); err != nil {
					log.Println("could not write to frame client, dropping frame and continuing in 250ms", err)

					time.Sleep(time.Millisecond * 250)

					continue
				}
			}
		}
	}(&wg)

	if *genesis {
		go func(wg *sync.WaitGroup) {
			log.Println("Reading from frame service")

			for {
				frame, err := frameService.Read()
				if err != nil {
					log.Println("could not read from frame service, dropping frame and continuing in 250ms", err)

					time.Sleep(time.Millisecond * 250)
				}

				if frame == nil {
					log.Println("read invalid frame from from frame service, dropping frame")

					continue
				}

				if *debug {
					log.Println("Read frame from from frame service")
				}

				if valid := preSharedKeyValidator.Validate(frame.PreSharedKey); !valid {
					log.Println("got invalid pre-shared key, dropping frame")

					continue
				}

				rawFrame, _, err := frameConverter.ToInternal(frame)
				if err != nil {
					log.Println("could not convert external frame to internal frame, dropping frame", err)

					continue
				}

				if *debug {
					log.Println("Writing frame from frame service to TAP device")
				}

				if err := tapDevice.Write(rawFrame); err != nil {
					log.Println("could not write to TAP device, dropping frame", err)

					continue
				}
			}
		}(&wg)
	} else {
		go func(wg *sync.WaitGroup) {
			log.Println("Reading from frame client")

			for {
				frame, err := frameClient.Read()
				if err != nil {
					log.Println("could not read from frame client, dropping frame and reconnecting in 250ms", err)

					time.Sleep(time.Millisecond * 250)

					if err := frameClient.Open(); err != nil {
						log.Println("could not reconnect to frame client, retrying in 250ms")
					}

					continue
				}

				if frame == nil {
					log.Println("read invalid frame from from frame client, dropping frame")

					continue
				}

				if valid := preSharedKeyValidator.Validate(frame.PreSharedKey); !valid {
					log.Println("got invalid pre-shared key, dropping frame")

					continue
				}

				rawFrame, _, err := frameConverter.ToInternal(frame)
				if err != nil {
					log.Println("could not convert external frame to internal frame, dropping frame", err)

					continue
				}

				if *debug {
					log.Println("Writing frame from frame client to TAP device")
				}

				if err := tapDevice.Write(rawFrame); err != nil {
					log.Println("could not write to TAP device, dropping frame", err)

					continue
				}
			}
		}(&wg)
	}

	wg.Wait()
}
