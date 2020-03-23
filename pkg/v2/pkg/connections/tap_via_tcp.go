package connections

import (
	"net"

	"github.com/pojntfx/gloeth/pkg/v2/pkg/constants"
)

// TAPviaTCPConnection is connection that writes from TAP to TCP and vice versa
type TAPviaTCPConnection struct {
	localAddr, remoteAddr *net.TCPAddr
	framesChan            chan []byte
}

// NewTAPviaTCPConnection creates a new TAP via TCP connection
func NewTAPviaTCPConnection(localAddr, remoteAddr *net.TCPAddr, framesChan chan []byte) *TAPviaTCPConnection {
	return &TAPviaTCPConnection{
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		framesChan: framesChan,
	}
}

// Open opens the TAP via TCP connection
func (t *TAPviaTCPConnection) Open() error {
	return nil
}

// Close closes the TAP via TCP connection
func (t *TAPviaTCPConnection) Close() error {
	return nil
}

// Write writes to the TAP via TCP connection
func (t *TAPviaTCPConnection) Write(frame []byte) (int, error) {
	conn, err := net.Dial("tcp", t.remoteAddr.String())
	if err != nil {
		return t.Write(frame) // Retry
	}

	defer conn.Close()
	return conn.Write(frame)
}

// Read reads from the TAP via TCP connection
func (t *TAPviaTCPConnection) Read() error {
	l, err := net.ListenTCP("tcp", t.localAddr)
	if err != nil {
		return err
	}

	for {
		frame := make([]byte, constants.TAP_FRAME_SIZE)

		conn, err := l.AcceptTCP()
		if err != nil {
			return err
		}

		defer conn.Close()
		_, err = conn.Read(frame)
		if err != nil {
			return err
		}

		t.framesChan <- frame
	}
}
