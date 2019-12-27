package pkg

import (
	"github.com/mdlayher/ethernet"
)

func GetDestinationMacAddressFromFrame(frame []byte) (error, string) {
	var unmarshalledFrame ethernet.Frame

	if err := (&unmarshalledFrame).UnmarshalBinary(frame); err != nil {
		return err, ""
	}

	destination := unmarshalledFrame.Destination.String()

	return nil, destination
}
