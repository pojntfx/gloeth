package main

import (
	"flag"
	"log"
	"net"

	"github.com/pojntfx/gloeth/pkg/proto/generated/proto"
	"github.com/pojntfx/gloeth/pkg/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	deviceName := flag.String("deviceName", "gloeth0", "Network device name")
	token := flag.String("preSharedKey", "supersecretkey", "Pre-shared key")

	localAddr := flag.String("localAddr", "0.0.0.0:1927", "Local address")
	localCert := flag.String("localCert", "local.crt", "Local certificate")
	localKey := flag.String("localKey", "local.key", "Local key")

	remoteAddr := flag.String("remoteAddr", "10.0.0.25:1927", "Remote address")
	remoteCert := flag.String("remoteCert", "remote.crt", "Remote certificate")

	debug := flag.Bool("debug", false, "Enable debugging ouput")

	flag.Parse()

	if *debug {
		log.Println(*deviceName, *token, *localAddr, *localCert, *localKey, *remoteAddr, *remoteCert, *debug)
	}

	listenAddress, err := net.ResolveTCPAddr("tcp", *localAddr)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.ListenTCP("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewServerTLSFromFile(*localCert, *localKey)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer(grpc.Creds(creds))

	reflection.Register(s)
	proto.RegisterFrameServiceServer(s, &services.FrameService{})

	if err = s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
