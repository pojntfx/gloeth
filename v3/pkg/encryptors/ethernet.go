package encryptors

// Ethernet encrypts and decrypts ethernet frames
type Ethernet struct {
	key string
}

// NewEthernet creates a new ethernet encryptor
func NewEthernet(key string) *Ethernet {
	return &Ethernet{key}
}

// Encrypt encrypts an ethernet frame
func (e *Ethernet) Encrypt(frame []byte) ([]byte, error) {
	return nil, nil
}

// Decrypt decrypts an ethernet frame
func (e *Ethernet) Decrypt(frame []byte) ([]byte, error) {
	return nil, nil
}
