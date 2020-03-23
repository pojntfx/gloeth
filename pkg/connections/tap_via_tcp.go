package connections

import (
	"log"
	"net"

	"github.com/pojntfx/gloeth/pkg/constants"
)

// TAPviaTCPConnection is connection that writes from TAP to TCP and vice versa
type TAPviaTCPConnection struct {
	localAddr, remoteAddr *net.TCPAddr
	framesChan            chan []byte
	conn                  *net.TCPConn
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
	conn, err := net.DialTCP("tcp", nil, t.remoteAddr)
	if err != nil {
		return t.Open()
	}

	t.conn = conn

	return nil
}

// Close closes the TAP via TCP connection
func (t *TAPviaTCPConnection) Close() error {
	return nil
}

// Write writes to the TAP via TCP connection
func (t *TAPviaTCPConnection) Write(frame []byte) (int, error) {
	log.Println("Writing frame")

	n, err := t.conn.Write(frame)
	if err != nil {
		// Retry
		if err := t.Close(); err != nil {
			return t.Write(frame)
		}

		if err := t.Open(); err != nil {
			return t.Write(frame)
		}

		return t.Write(frame)
	}

	log.Println("Wrote frame")

	return n, nil
}

// Read reads from the TAP via TCP connection
func (t *TAPviaTCPConnection) Read() error {
	l, err := net.ListenTCP("tcp", t.localAddr)
	if err != nil {
		return err
	}

	for {
		frame := make([]byte, constants.FRAME_SIZE)

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
