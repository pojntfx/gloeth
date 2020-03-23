package constants

import "time"

const (
	TAP_FRAME_SIZE        = 1472
	TCP_FRAME_SIZE        = TAP_FRAME_SIZE + 14
	TIMESTAMP_SIZE        = 8
	FRAME_DISCARD_TIMEOUT = time.Second * 3
)
