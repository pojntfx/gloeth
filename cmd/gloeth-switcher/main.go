package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/pojntfx/ethernet"
	"github.com/pojntfx/gloeth/pkg/constants"
	"github.com/pojntfx/gloeth/pkg/switchers"
)

func main() {
	laddrFlag := flag.String("laddr", ":1234", "Listen address")
	flag.Parse()

	laddr, err := net.ResolveTCPAddr("tcp", *laddrFlag)
	if err != nil {
		log.Fatal(err)
	}

	frameChan := make(chan ethernet.Frame)

	switcher := switchers.NewMACviaTCPSwitcher(laddr, constants.FRAME_SIZE, constants.TIMESTAMP_SIZE, frameChan)

	defer switcher.Close()
	if err := switcher.Open(); err != nil {
		log.Fatal(err)
	}

	go switcher.Read()

	for {
		frame := <-frameChan

		log.Println(frame.Source, frame.Destination, string(frame.Payload))

		go func() {
			var err error

			for i := 0; i < 10; i++ {
				err = switcher.Write(frame)
				if err == nil {
					return
				}

				log.Println(err)

				time.Sleep(time.Millisecond * 50)
			}

			log.Printf("frame discarded, could not write: %v\n", err)
		}()

	}
}
