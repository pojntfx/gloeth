package main

import (
	"flag"
	"log"
)

func main() {
	deviceName := *flag.String("deviceName", "gloeth0", "Network device name")

	localAddr := *flag.String("localAddr", "0.0.0.0:1927", "Local address")
	localCert := *flag.String("localCert", "local.crt", "Local certificate")
	localKey := *flag.String("localKey", "local.key", "Local key")
	localToken := *flag.String("localToken", "supersecrettoken", "Local token")

	remoteAddr := *flag.String("remoteAddr", "10.0.0.25:1927", "Remote address")
	remoteCert := *flag.String("remoteCert", "remote.crt", "Remote certificate")
	remoteToken := *flag.String("remoteToken", "supersecrettoken", "Remote token")

	flag.Parse()

	log.Println(deviceName, localAddr, localCert, localKey, localToken, remoteAddr, remoteCert, remoteToken)
}
