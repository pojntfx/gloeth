package utils

import "net"

// SwitchConnection is a connection to a switch
type SwitchConnection struct {
	remoteAddr *net.TCPAddr
}

// NewSwitchConnection creates a new switch connection
func NewSwitchConnection(remoteAddr *net.TCPAddr) *SwitchConnection {
	return &SwitchConnection{remoteAddr: remoteAddr}
}

// Write writes a packet to the switch
func (s *SwitchConnection) Write(frame []byte) error {
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

// Read reads a packet from the switch
func (s *SwitchConnection) Read(packetChan chan []byte) error {
	return nil
}
