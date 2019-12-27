package pkg

import (
	"bytes"
	"encoding/base64"
	"github.com/mdlayher/ethernet"
	"strings"
)

func GetDestinationMacAddressFromFrame(frame []byte) (error, string) {
	var unmarshalledFrame ethernet.Frame

	if err := (&unmarshalledFrame).UnmarshalBinary(frame); err != nil {
		return err, ""
	}

	destination := unmarshalledFrame.Destination.String()

	return nil, destination
}

func Encode(fromTAP []byte) (error, string) {
	toTCP := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, toTCP)

	encoder.Write(fromTAP)

	encoder.Close()

	return nil, toTCP.String() + "\n"
}

func Decode(fromTCP string) (error, []byte) {
	fromTCPWithoutNewline := strings.TrimSuffix(fromTCP, "\n")

	toTAP, err := base64.StdEncoding.DecodeString(fromTCPWithoutNewline)

	return err, toTAP
}
