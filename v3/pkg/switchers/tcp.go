package switchers

import (
	"log"
	"net"

	"github.com/pojntfx/gloeth/v3/pkg/wrappers"
)

// TCP switches TCP connections
type TCP struct {
	readChan chan [wrappers.WrappedFrameSize]byte
	laddr    *net.TCPAddr
	listener *net.TCPListener
	conns    []map[string]*net.TCPConn
}

// NewTCP creates a new TCP switcher
func NewTCP(readChan chan [wrappers.WrappedFrameSize]byte, laddr *net.TCPAddr) *TCP {
	return &TCP{readChan, laddr, nil, nil}
}

// Open opens the TCP switcher
func (t *TCP) Open() error {
	l, err := net.ListenTCP("tcp", t.laddr)
	if err != nil {
		return err
	}

	t.listener = l

	return nil
}

// Close closes the TCP switcher
func (t *TCP) Close() error {
	return t.listener.Close()
}

// Read reads from the TCP switcher
func (t *TCP) Read() error {
	for {
		conn, err := t.listener.AcceptTCP()
		if err != nil {
			return err
		}

		// TODO: Decompose to main method
		go func() {
			for {
				readFrame := [wrappers.WrappedFrameSize]byte{}

				_, err := conn.Read(readFrame[:])
				if err != nil {
					log.Fatal(err)
				}

				t.readChan <- readFrame
			}
		}()
	}
}

// Write writes to a connection on the TCP switcher
func (t *TCP) Write(conn *net.TCPConn, frame [wrappers.WrappedFrameSize]byte) error {
	return nil
}

// GetConnectionsForMAC gets the connections for a given MAC address
func (t *TCP) GetConnectionsForMAC(mac *net.HardwareAddr) ([]*net.TCPConn, error) {
	return []*net.TCPConn{}, nil
}
