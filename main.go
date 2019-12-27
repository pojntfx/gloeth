package main

import (
	"flag"
	redisLib "github.com/go-redis/redis/v7"
	"github.com/pojntfx/gloeth/pkg"
	"log"
)

func main() {
	var (
		tapName         = flag.String("device", "gloeth0", "Name of the network device to create")
		tcpReadHostPort = flag.String("listen", "127.0.0.1:1234", "Host:port to listen on")

		redisHostPort = flag.String("redis-host", "127.0.0.1:6379", "Host:port of Redis")
		redisPassword = flag.String("redis-password", "", "Password for Redis")
	)
	flag.Parse()

	var (
		tapErrorChan = make(chan error)
		tcpErrorChan = make(chan error)

		tapStatusChan = make(chan string)
		tcpStatusChan = make(chan string)

		tapReadFramesChan = make(chan []byte)
		tcpReadFramesChan = make(chan []byte)
	)

	redis := pkg.Redis{}

	redis.Connect(*redisHostPort, *redisPassword)

	tap := pkg.TAP{
		Name: *tapName,
	}

	if err := tap.Init(); err != nil {
		log.Fatalln("TAP init error:", err)
	}

	err, macAddress := tap.GetMacAddress()
	if err != nil {
		log.Fatalln("TAP registration error:", err)
	}

	if err := redis.RegisterNode(macAddress, *tcpReadHostPort); err != nil {
		log.Fatalln("TAP registration error:", err)
	}

	tcp := pkg.TCP{
		ReadHostPort: *tcpReadHostPort,
	}

	go tap.Read(tapErrorChan, tapStatusChan, tapReadFramesChan)
	go tcp.Read(tcpErrorChan, tcpStatusChan, tcpReadFramesChan)

	for {
		select {
		case err := <-tapErrorChan:
			log.Println("TAP error:", err)
		case err := <-tcpErrorChan:
			log.Println("TCP error:", err)

		case status := <-tapStatusChan:
			log.Println("TAP status:", status)
		case status := <-tcpStatusChan:
			log.Println("TCP status:", status)

		case frame := <-tapReadFramesChan:
			err, macAddress := pkg.GetDestinationMacAddressFromFrame(frame)
			if err != nil {
				log.Println("SWITCH error:", err)
			}

			destinations, err := redis.GetTcpReadHostPortsForMacAddress(macAddress)
			if err == redisLib.Nil {
				log.Println("SWITCH error: could not find node with mac address", macAddress, err)
			} else if err != nil {
				log.Println("SWITCH error:", err)
			} else {
				for _, destination := range destinations {
					go tcp.Write(tcpErrorChan, tcpStatusChan, frame, destination)
				}
			}
		case frame := <-tcpReadFramesChan:
			go tap.Write(tapErrorChan, tapStatusChan, frame)
		}
	}
}
