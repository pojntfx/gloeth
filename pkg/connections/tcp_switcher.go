package connections

import (
	"net"

	"github.com/pojntfx/gloeth/pkg/wrappers"
)

// TCPSwitcher is a connection to a TCP switcher
type TCPSwitcher struct {
	readChan chan [wrappers.WrappedFrameSize]byte
	raddr    *net.TCPAddr
	conn     *net.TCPConn
}

// NewTCPSwitcher creates a new TCP switcher connection
func NewTCPSwitcher(readChan chan [wrappers.WrappedFrameSize]byte, raddr *net.TCPAddr) *TCPSwitcher {
	return &TCPSwitcher{readChan, raddr, nil}
}

// Open opens the connection to the TCP switcher
func (t *TCPSwitcher) Open() error {
	conn, err := net.DialTCP("tcp", nil, t.raddr)
	if err != nil {
		return err
	}

	t.conn = conn

	return nil
}

// Close closes the connection to the TCP switcher
func (t *TCPSwitcher) Close() error {
	return t.conn.Close()
}

// Read reads from the TCP switcher
func (t *TCPSwitcher) Read() error {
	for {
		readFrame := [wrappers.WrappedFrameSize]byte{}

		_, err := t.conn.Read(readFrame[:])
		if err != nil {
			return err
		}

		t.readChan <- readFrame
	}
}

// Write writes to the TCP switcher
func (t *TCPSwitcher) Write(frame [wrappers.WrappedFrameSize]byte) error {
	_, err := t.conn.Write(frame[:])

	return err
}
