package connections

import "net"

// TCPSwitcher is a connection to a TCP switcher
type TCPSwitcher struct {
	readChan   chan []byte
	remoteAddr *net.TCPAddr
}

// NewTCPSwitcher creates a new TCP switcher connection
func NewTCPSwitcher(readChan chan []byte, remoteAddr *net.TCPAddr) *TCPSwitcher {
	return &TCPSwitcher{readChan, remoteAddr}
}

// Open opens the connection to the TCP switcher
func (t *TCPSwitcher) Open() error {
	return nil
}

// Close closes the connection to the TCP switcher
func (t *TCPSwitcher) Close() error {
	return nil
}

// Read reads from the TCP switcher
func (t *TCPSwitcher) Read() error {
	return nil
}

// Write writes to the TCP switcher
func (t *TCPSwitcher) Write(frame []byte) error {
	return nil
}
