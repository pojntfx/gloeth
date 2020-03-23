package protocol

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/pojntfx/gloeth/pkg/v2/pkg/constants"
)

// Encoder encodes and decodes frames
type Encoder struct {
}

// NewEncoder creates a new encoder
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encapsulate encapsulates a frame
// A frame is composed of a nanosecond timestamp (8 bytes) and a plaintext frame (1-1464 bytes)
func (e *Encoder) Encapsulate(frame []byte) ([]byte, error) {
	timeInByte := make([]byte, constants.TIMESTAMP_SIZE)
	binary.BigEndian.PutUint64(timeInByte, uint64(time.Now().UnixNano()))

	return append(timeInByte, frame...), nil
}

// Decapsulate decapsulates a frame
func (e *Encoder) Decapsulate(frame []byte) ([]byte, error) {
	if len(frame) < constants.TIMESTAMP_SIZE+1 {
		return nil, errors.New("invalid frame size")
	}

	timeInByte := frame[0:constants.TIMESTAMP_SIZE]
	timeInUnixNano := int64(binary.BigEndian.Uint64(timeInByte))
	if (time.Now().UnixNano() - timeInUnixNano) > constants.FRAME_DISCARD_TIMEOUT.Nanoseconds() {
		return nil, errors.New("timestamp out of acceptable range")
	}

	return frame[constants.TIMESTAMP_SIZE:], nil
}
