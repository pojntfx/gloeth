package switchers

import (
	"fmt"
	"net"

	gm "github.com/cseeger-epages/mac-gen-go"
)

const (
	SwitcherInfoSize = 17 // Size of a switcher info
)

// SwitcherInfo provides information about a switcher
type SwitcherInfo struct {
	laddr    *net.TCPAddr
	listener *net.TCPListener
	mac      *net.HardwareAddr
}

// NewSwitcherInfo creates a new SwitcherInfo
func NewSwitcherInfo(laddr *net.TCPAddr) *SwitcherInfo {
	return &SwitcherInfo{
		laddr,
		nil,
		nil,
	}
}

// RequestMACAddress assigns a MAC address to the switcher
func (s *SwitcherInfo) RequestMACAddress() error {
	prefix := gm.GenerateRandomLocalMacPrefix(false)

	suffix, err := gm.CalculateNICSufix(net.ParseIP("10.0.0.1"))
	if err != nil {
		return err
	}

	rawMAC := fmt.Sprintf("%v:%v", prefix, suffix)

	mac, err := net.ParseMAC(rawMAC)
	if err != nil {
		return err
	}

	s.mac = &mac

	return nil
}

// Open opens the switcher info
func (s *SwitcherInfo) Open() error {
	l, err := net.ListenTCP("tcp", s.laddr)
	if err != nil {
		return err
	}

	s.listener = l

	return nil
}

// Close closes the switcher info
func (s *SwitcherInfo) Close() error {
	return s.listener.Close()
}

// Read reads from the switcher info
func (s *SwitcherInfo) Read() error {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			return err
		}

		outMAC := [SwitcherInfoSize]byte{}
		copy(outMAC[:], s.mac.String())

		if _, err := conn.Write(outMAC[:]); err != nil {
			return err
		}
	}
}
