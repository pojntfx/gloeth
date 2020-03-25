package encryptors

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"github.com/pojntfx/gloeth/pkg/wrappers"
)

const (
	PlaintextFrameSize = wrappers.EncryptedFrameSize - aes.BlockSize // PlaintextFrameSize is the size of a decrypted frame
)

// Ethernet encrypts and decrypts ethernet frames
type Ethernet struct {
	key string
}

// NewEthernet creates a new ethernet encryptor
func NewEthernet(key string) *Ethernet {
	return &Ethernet{key}
}

// Encrypt encrypts an ethernet frame
func (e *Ethernet) Encrypt(frame [PlaintextFrameSize]byte) ([wrappers.EncryptedFrameSize]byte, error) {
	block, err := aes.NewCipher([]byte(e.key))
	if err != nil {
		return [wrappers.EncryptedFrameSize]byte{}, err
	}

	encryptedFrame := [wrappers.EncryptedFrameSize]byte{}
	iv := encryptedFrame[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return [wrappers.EncryptedFrameSize]byte{}, nil
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	stream.XORKeyStream(encryptedFrame[aes.BlockSize:], frame[:])

	return encryptedFrame, nil
}

// Decrypt decrypts an ethernet frame
func (e *Ethernet) Decrypt(frame [wrappers.EncryptedFrameSize]byte) ([PlaintextFrameSize]byte, error) {
	block, err := aes.NewCipher([]byte(e.key))
	if err != nil {
		return [PlaintextFrameSize]byte{}, err
	}

	iv := frame[:aes.BlockSize]
	encryptedFrame := frame[aes.BlockSize:]
	decryptedFrame := [PlaintextFrameSize]byte{}

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(decryptedFrame[:], encryptedFrame)

	return decryptedFrame, nil
}
