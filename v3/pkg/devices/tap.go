package devices

// TAP is a TAP device
type TAP struct {
	readChan chan []byte
	name     string
}

// NewTAP creates a new TAP device
func NewTAP(readChan chan []byte, name string) *TAP {
	return &TAP{readChan, name}
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
func (t *TAP) Write(frame []byte) error {
	return nil
}
