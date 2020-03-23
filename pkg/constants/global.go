package constants

import "time"

const (
	TIMESTAMP_SIZE    = 8               // TIMESTAMP_SIZE is the size of the timestamp
	TIMESTAMP_TIMEOUT = time.Second * 3 // TIMESTAMP_TIMEOUT is the maximum duration after which a frame will be discarded

	TCP_MTU = 1472                          // The TCP MTU
	TAP_MTU = TCP_MTU - 14 - TIMESTAMP_SIZE // TAP_MTU is TCP_MTU - ethernet header (14) - TIMESTAMP_SIZE
)
