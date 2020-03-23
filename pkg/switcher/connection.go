package switcher

import "net"

// Connection is a connection to a switcher
type Connection struct {
	remoteAddr *net.TCPAddr
}

// NewConnection creates a new switcher connection
func NewConnection(remoteAddr *net.TCPAddr) *Connection {
	return &Connection{remoteAddr: remoteAddr}
}

// Write writes a packet to the switcher
func (s *Connection) Write(frame []byte) error {
	conn, err := net.Dial("tcp", s.remoteAddr.String())
	if err != nil {
		return err
	}

	_, err = conn.Write(frame)
	if err != nil {
		return err
	}

	return conn.Close()
}

// Read reads a packet from the switcher
func (s *Connection) Read(packetChan chan []byte) error {
	return nil
}
