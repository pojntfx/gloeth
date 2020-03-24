package devices

import "github.com/pojntfx/gloeth/v3/pkg/encryptors"

const (
	MTU = encryptors.PlaintextFrameSize - 14 // MTU is the MTU, which is the plaintext frame size - ethernet header (14 bytes)
)

// TAP is a TAP device
type TAP struct {
	readChan chan [encryptors.PlaintextFrameSize]byte
	mtu      uint
	name     string
}

// NewTAP creates a new TAP device
func NewTAP(readChan chan [encryptors.PlaintextFrameSize]byte, mtu uint, name string) *TAP {
	return &TAP{readChan, mtu, name}
}

// Open opens the TAP device
func (t *TAP) Open() error {
	return nil
}

// Close closes the TAP device
func (t *TAP) Close() error {
	return nil
}

// Read reads from the TAP device
func (t *TAP) Read() error {
	return nil
}

// Write writes from the TAP device
func (t *TAP) Write(frame [encryptors.PlaintextFrameSize]byte) error {
	return nil
}
