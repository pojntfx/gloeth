package connections

import (
	"net"

	"github.com/pojntfx/gloeth/pkg/switchers"
)

// SwitcherInfo is a connection to a switcher info
type SwitcherInfo struct {
	readChan chan *net.HardwareAddr
	raddr    *net.TCPAddr
}

// NewSwitcherInfo creates a new switcher info connection
func NewSwitcherInfo(readChan chan *net.HardwareAddr, raddr *net.TCPAddr) *SwitcherInfo {
	return &SwitcherInfo{readChan, raddr}
}

// Read reads from the switcher info
func (t *SwitcherInfo) Read() error {
	for {
		conn, err := GetConn(t.raddr)
		if err != nil {
			return err
		}

		macRaw := [switchers.SwitcherInfoSize]byte{}

		if _, err := conn.Read(macRaw[:]); err != nil {
			return err
		}

		mac, err := net.ParseMAC(string(macRaw[:]))
		if err != nil {
			return err
		}

		t.readChan <- &mac
	}
}
