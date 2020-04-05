package switchers

import (
	"fmt"
	"net"

	cmap "github.com/orcaman/concurrent-map"
	"github.com/pojntfx/gloeth/pkg/wrappers"
)

// TCP switches TCP connections
type TCP struct {
	readChan chan [wrappers.WrappedFrameSize]byte
	connChan chan *net.TCPConn
	laddr    *net.TCPAddr
	listener *net.TCPListener
	conns    cmap.ConcurrentMap
}

// NewTCP creates a new TCP switcher
func NewTCP(readChan chan [wrappers.WrappedFrameSize]byte, connChan chan *net.TCPConn, laddr *net.TCPAddr, conns map[string]*net.TCPConn) *TCP {
	iconns := cmap.New()

	for mac, conn := range conns {
		iconns.Set(mac, conn)
	}

	return &TCP{readChan, connChan, laddr, nil, iconns}
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
	t.conns.Set(mac.String(), conn)
}

// GetConnectionsForMAC gets the connections for a given MAC address
func (t *TCP) GetConnectionsForMAC(destMAC, srcMAC *net.HardwareAddr) ([]*net.TCPConn, error) {
	dest := destMAC.String()
	src := srcMAC.String()

	if dest == "ff:ff:ff:ff:ff:ff" {
		connsToReturn := []*net.TCPConn{}

		for connt := range t.conns.Iter() {
			if connt.Key != src {
				connsToReturn = append(connsToReturn, connt.Val.(*net.TCPConn))
			}
		}

		return connsToReturn, nil
	}

	conn, ok := t.conns.Get(dest)
	if !ok {
		return []*net.TCPConn{}, fmt.Errorf("no connection found for dest %v", dest)
	}

	return []*net.TCPConn{conn.(*net.TCPConn)}, nil
}

// Write writes to a connection on the TCP switcher
func (t *TCP) Write(conn *net.TCPConn, frame [wrappers.WrappedFrameSize]byte) error {
	_, err := conn.Write(frame[:])

	return err
}
