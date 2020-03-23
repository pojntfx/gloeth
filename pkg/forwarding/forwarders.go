package forwarding

import (
	"log"
	"net"

	"github.com/pojntfx/gloeth/pkg/constants"
	"github.com/pojntfx/gloeth/pkg/encoding"
	"github.com/pojntfx/gloeth/pkg/switcher"
	"github.com/pojntfx/gloeth/pkg/tap"
)

// ForwardTCPtoTAP forwards TCP packets to a TAP device
func ForwardTCPtoTAP(tcpConn *net.TCPListener, tapDevice *tap.Device, remoteAddr *net.TCPAddr) {
	log.Printf("forwarding TCP to TAP with remote %v:%v\n", remoteAddr.IP, remoteAddr.Port)

	for {
		packet := make([]byte, constants.TAP_MTU)
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

		frame, invalid := encoding.DecapsulateFrame(packet[0:n])
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
func ForwardTAPtoTCP(tapDevice *tap.Device, remoteAddr *net.TCPAddr) {
	log.Printf("forwarding TAP to TCP with remote %v:%v\n", remoteAddr.IP, remoteAddr.Port)

	for {
		frame := make([]byte, constants.TAP_MTU+14)
		var encFrame []byte

		n, err := tapDevice.Read(frame)
		if err != nil {
			log.Fatalf("could not read from TAP device: %v\n", err)
		}

		encFrame, invalid := encoding.EncapsulateFrame(frame[0:n])
		if invalid != nil {
			continue
		}

		s := switcher.NewConnection(remoteAddr)

		if err := s.Write(encFrame); err != nil {
			log.Printf("could not dial %v, retrying", remoteAddr)

			ForwardTAPtoTCP(tapDevice, remoteAddr)
		}
	}
}
