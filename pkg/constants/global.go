package constants

import "time"

const (
	FRAME_SIZE            = 1472                             // FRAME_SIZE is the maximum frame size
	TIMESTAMP_SIZE        = 8                                // TIMESTAMP_SIZE is the size of the timestamp header
	MTU                   = FRAME_SIZE - 14 - TIMESTAMP_SIZE // MTU is Frame size - Ethernet Header (not included in the MTU) - Timestamp Size
	FRAME_DISCARD_TIMEOUT = time.Second * 3                  // FRAME_DISCARD_TIMEOUT is the duration after which to discard a packet
)
