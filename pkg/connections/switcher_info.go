package connections

import (
	"io/ioutil"
	"net"
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

		macRaw, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}

		mac, err := net.ParseMAC(string(macRaw))
		if err != nil {
			return err
		}

		t.readChan <- &mac
	}
}
