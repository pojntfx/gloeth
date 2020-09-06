package main

import (
	"flag"
	"log"
)

func main() {
	deviceName := *flag.String("deviceName", "gloeth0", "Network device name")
	token := *flag.String("preSharedKey", "supersecretkey", "Pre-shared key")

	localAddr := *flag.String("localAddr", "0.0.0.0:1927", "Local address")
	localCert := *flag.String("localCert", "local.crt", "Local certificate")
	localKey := *flag.String("localKey", "local.key", "Local key")

	remoteAddr := *flag.String("remoteAddr", "10.0.0.25:1927", "Remote address")
	remoteCert := *flag.String("remoteCert", "remote.crt", "Remote certificate")

	flag.Parse()

	log.Println(deviceName, token, localAddr, localCert, localKey, remoteAddr, remoteCert)
}
