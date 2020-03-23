package switchers

import "net"

// TCP switches TCP connections
type TCP struct {
	readChan   chan []byte
	listenAddr *net.TCPAddr
}

// NewTCP creates a new TCP switcher
func NewTCP(readChan chan []byte, listenAddr *net.TCPAddr) *TCP {
	return &TCP{readChan, listenAddr}
}

// Open opens the TCP switcher
func (t *TCP) Open() error {
	return nil
}

// Close closes the TCP switcher
func (t *TCP) Close() error {
	return nil
}

// Read reads from the TCP switcher
func (t *TCP) Read() error {
	return nil
}

// Write writes to a connection on the TCP switcher
func (t *TCP) Write(conn *net.TCPConn, frame []byte) error {
	return nil
}

// GetConnectionsForMAC gets the connections for a given MAC address
func (t *TCP) GetConnectionsForMAC(mac *net.HardwareAddr) ([]*net.TCPConn, error) {
	return []*net.TCPConn{}, nil
}
