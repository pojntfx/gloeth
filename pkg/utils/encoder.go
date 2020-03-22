package utils

import (
	"encoding/binary"
	"errors"
	"time"
)

// EncapsulateFrame encapsulates a frame
// A frame is composed of a nanosecond timestamp (8 bytes) and a plaintext frame (1-1432 bytes)
func EncapsulateFrame(frame []byte) ([]byte, error) {
	timeInByte := make([]byte, 8)
	binary.BigEndian.PutUint64(timeInByte, uint64(time.Now().UnixNano()))

	return append(timeInByte, frame...), nil
}

// DecapsulateFrame decapsulates a frame
func DecapsulateFrame(frame []byte) ([]byte, error) {
	if len(frame) < (TIMESTAMP_SIZE + 1) {
		return nil, errors.New("invalid encapsulated frame size")
	}

	timeInByte := frame[0:TIMESTAMP_SIZE]
	timeInUnixNano := int64(binary.BigEndian.Uint64(timeInByte))
	if (time.Now().UnixNano() - timeInUnixNano) > TIMESTAMP_TIMEOUT.Nanoseconds() {
		return nil, errors.New("timestamp out of acceptable range")
	}

	return frame[TIMESTAMP_SIZE:], nil
}
