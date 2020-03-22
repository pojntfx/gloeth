package utils

import (
	"log"
	"net"
	"time"
)

const (
	TIMESTAMP_SIZE    = 8               // TIMESTAMP_SIZE is the size of the timestamp
	TIMESTAMP_TIMEOUT = time.Second * 3 // TIMESTAMP_TIMEOUT is the maximum duration after which a frame will be discarded

	TCP_MTU = 1472                          // The TCP MTU
	TAP_MTU = TCP_MTU - 14 - TIMESTAMP_SIZE // TAP_MTU is TCP_MTU - ethernet header (14) - TIMESTAMP_SIZE
)

// ForwardTCPtoTAP forwards TCP packets to a TAP device
func ForwardTCPtoTAP(tcpConn *net.TCPListener, tapDevice *TAPDevice, remoteAddr *net.TCPAddr) {
	log.Printf("forwarding TCP to TAP with remote %v:%v\n", remoteAddr.IP, remoteAddr.Port)

	for {
		packet := make([]byte, TCP_MTU)
		var frame []byte

		conn, err := tcpConn.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}

		n, err := conn.Read(packet)
		if err != nil {
			log.Fatalf("could not read from TCP socket: %v\n", err)
		}

		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}

		frame, invalid := DecapsulateFrame(packet[0:n])
		if invalid != nil {
			continue
		}

		_, err = tapDevice.Write(frame)
		if err != nil {
			log.Fatalf("could not write to TAP device: %v\n", err)
		}
	}
}

// ForwardTAPtoTCP forwards frames from a TAP device to a TCP connection
func ForwardTAPtoTCP(tapDevice *TAPDevice, remoteAddr *net.TCPAddr) {
	log.Printf("forwarding TAP to TCP with remote %v:%v\n", remoteAddr.IP, remoteAddr.Port)

	for {
		frame := make([]byte, TAP_MTU+14)
		var encFrame []byte

		n, err := tapDevice.Read(frame)
		if err != nil {
			log.Fatalf("could not read from TAP device: %v\n", err)
		}

		encFrame, invalid := EncapsulateFrame(frame[0:n])
		if invalid != nil {
			continue
		}

		conn, err := net.Dial("tcp", remoteAddr.String())
		if err != nil {
			log.Printf("could not dial %v, retrying", remoteAddr)

			ForwardTAPtoTCP(tapDevice, remoteAddr)

			continue
		}

		_, err = conn.Write(encFrame)
		if err != nil {
			log.Fatalf("could not write to TCP socket: %v\n", err)
		}

		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
	}
}
