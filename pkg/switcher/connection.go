package switcher

import (
	"fmt"
	"net"

	"github.com/pojntfx/gloeth/pkg/constants"
)

// ReadPacket is a packet that has been read from the switcher
type ReadPacket struct {
	Length  int
	Payload []byte
}

// Connection is a connection to a switcher
type Connection struct {
	localAddr, remoteAddr *net.TCPAddr
}

// NewConnection creates a new switcher connection
func NewConnection(localAddr, remoteAddr *net.TCPAddr) *Connection {
	return &Connection{localAddr: localAddr, remoteAddr: remoteAddr}
}

// Write writes a packet to the switcher
func (s *Connection) Write(frame []byte) error {
	conn, err := net.Dial("tcp", s.remoteAddr.String())
	if err != nil {
		return fmt.Errorf("could not dial TCP: %v", err)
	}

	_, err = conn.Write(frame)
	if err != nil {
		return fmt.Errorf("could not write to TCP connection: %v", err)
	}

	return conn.Close()
}

// Read reads a packet from the switcher
func (s *Connection) Read(packetChan chan ReadPacket) error {
	tcpListener, err := net.ListenTCP("tcp", s.localAddr)
	if err != nil {
		return err
	}

	for {
		packet := make([]byte, constants.TAP_MTU)

		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			return err
		}

		n, err := conn.Read(packet)
		if err != nil {
			return fmt.Errorf("could not read from TCP connection: %v", err)
		}

		packetChan <- ReadPacket{n, packet}

		if err := conn.Close(); err != nil {
			return fmt.Errorf("could not close TCP connection: %v", err)
		}
	}
}
