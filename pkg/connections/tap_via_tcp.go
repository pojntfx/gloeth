package connections

import (
	"log"
	"net"
)

// TAPviaTCPConnection is connection that writes from TAP to TCP and vice versa
type TAPviaTCPConnection struct {
	remoteAddr *net.TCPAddr
	framesChan chan []byte
	conn       *net.TCPConn
	frameSize  uint
}

// NewTAPviaTCPConnection creates a new TAP via TCP connection
func NewTAPviaTCPConnection(remoteAddr *net.TCPAddr, frameSize uint, framesChan chan []byte) *TAPviaTCPConnection {
	return &TAPviaTCPConnection{
		remoteAddr: remoteAddr,
		framesChan: framesChan,
		frameSize:  frameSize,
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
	for {
		frame := make([]byte, t.frameSize)
		if _, err := t.conn.Read(frame); err != nil {
			log.Fatal(err)
		}

		t.framesChan <- frame
	}
}
