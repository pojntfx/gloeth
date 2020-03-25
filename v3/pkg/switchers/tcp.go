package switchers

import (
	"fmt"
	"net"

	"github.com/pojntfx/gloeth/v3/pkg/wrappers"
)

// TCP switches TCP connections
type TCP struct {
	readChan chan [wrappers.WrappedFrameSize]byte
	connChan chan *net.TCPConn
	laddr    *net.TCPAddr
	listener *net.TCPListener
	conns    map[string]*net.TCPConn
}

// NewTCP creates a new TCP switcher
func NewTCP(readChan chan [wrappers.WrappedFrameSize]byte, connChan chan *net.TCPConn, laddr *net.TCPAddr) *TCP {
	return &TCP{readChan, connChan, laddr, nil, nil}
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

		t.connChan <- conn
	}
}

// HandleFrame handles a frame
func (t *TCP) HandleFrame(frame [wrappers.WrappedFrameSize]byte) {
	t.readChan <- frame
}

// Register registers a TCP connection
func (t *TCP) Register(mac *net.HardwareAddr, conn *net.TCPConn) {
	t.conns[mac.String()] = conn
}

// GetConnectionsForMAC gets the connections for a given MAC address
func (t *TCP) GetConnectionsForMAC(destMAC, srcMAC *net.HardwareAddr) ([]*net.TCPConn, error) {
	dest := destMAC.String()
	src := srcMAC.String()

	if dest == "ff:ff:ff:ff:ff:ff" {
		connsToReturn := []*net.TCPConn{}

		for connDest, conn := range t.conns {
			if connDest != src {
				connsToReturn = append(connsToReturn, conn)
			}
		}

		return connsToReturn, nil
	}

	conn := t.conns[dest]
	if conn == nil {
		return []*net.TCPConn{}, fmt.Errorf("no connection found for dest %v", dest)
	}

	return []*net.TCPConn{conn}, nil
}

// Write writes to a connection on the TCP switcher
func (t *TCP) Write(conn *net.TCPConn, frame [wrappers.WrappedFrameSize]byte) error {
	return nil
}
