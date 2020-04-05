package switchers

import (
	"fmt"
	"net"

	gm "github.com/cseeger-epages/mac-gen-go"
)

// SwitcherInfo provides information about a switcher
type SwitcherInfo struct {
	mac *net.HardwareAddr
}

// NewSwitcherInfo creates a new SwitcherInfo
func NewSwitcherInfo() *SwitcherInfo {
	return &SwitcherInfo{}
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
