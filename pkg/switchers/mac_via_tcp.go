package switchers

import (
	"fmt"
	"log"
	"net"

	"github.com/pojntfx/ethernet"
	"github.com/pojntfx/gloeth/pkg/protocol"
)

// MACviaTCPSwitcher maps MAC addresses to TCP connections
type MACviaTCPSwitcher struct {
	listener              *net.TCPListener
	conns                 map[string]*net.TCPConn
	laddr                 *net.TCPAddr
	frameSize, headerSize uint
	frameChan             chan ethernet.Frame
}

// NewMACviaTCPSwitcher creates a new MAC via TCP switcher
func NewMACviaTCPSwitcher(laddr *net.TCPAddr, frameSize, headerSize uint, frameChan chan ethernet.Frame) *MACviaTCPSwitcher {
	return &MACviaTCPSwitcher{
		conns:      make(map[string]*net.TCPConn),
		laddr:      laddr,
		frameSize:  frameSize,
		headerSize: headerSize,
		frameChan:  frameChan,
	}
}

// Open opens the MAC via TCP switcher
func (m *MACviaTCPSwitcher) Open() error {
	l, err := net.ListenTCP("tcp", m.laddr)
	if err != nil {
		return err
	}

	m.listener = l

	return nil
}

// Close closes the MAC via TCP switcher
func (m *MACviaTCPSwitcher) Close() error {
	for _, conn := range m.conns {
		if err := conn.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Write writes a frame to a MAC address over TCP
func (m *MACviaTCPSwitcher) Write(frame ethernet.Frame) error {
	mac := frame.Destination.String()
	enc := protocol.NewEncoder()

	var connsToSendTo []*net.TCPConn
	if mac == "ff:ff:ff:ff:ff:ff" {
		for mac, conn := range m.conns {
			if mac != frame.Source.String() { // Don't send broadcast back to sender
				connsToSendTo = append(connsToSendTo, conn)
			}
		}
	} else {
		conn := m.conns[mac]
		connsToSendTo = append(connsToSendTo, conn)
		if conn == nil {
			return fmt.Errorf("no connection for MAC %v found", mac)
		}
	}

	binFrame, err := frame.MarshalBinary()
	if err != nil {
		return err
	}

	outFrame, err := enc.Encapsulate(binFrame)
	if err != nil {
		return err
	}

	for _, conn := range connsToSendTo {
		if _, err := conn.Write(outFrame); err != nil {
			return err
		}
	}

	return nil
}

// Read reads frames from the switcher
func (m *MACviaTCPSwitcher) Read() error {
	enc := protocol.NewEncoder()

	for {
		conn, err := m.listener.AcceptTCP()
		if err != nil {
			return err
		}

		go func() {
			for {
				inFrame := make([]byte, m.frameSize)
				if _, err := conn.Read(inFrame); err != nil {
					log.Fatal(err)
				}

				decFrame, err := enc.Decapsulate(inFrame)
				if err != nil {
					continue
				}

				var ethFrame ethernet.Frame
				ethFrame.UnmarshalBinary(decFrame)

				m.conns[ethFrame.Source.String()] = conn

				m.frameChan <- ethFrame
			}
		}()
	}
}
