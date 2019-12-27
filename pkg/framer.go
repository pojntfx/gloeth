package pkg

import (
	"bytes"
	"encoding/base64"
	"github.com/mdlayher/ethernet"
	"strings"
)

func GetDestinationMacAddressFromFrame(frame []byte) (error, string) {
	var unmarshalFrame ethernet.Frame

	if err := (&unmarshalFrame).UnmarshalBinary(frame); err != nil {
		return err, ""
	}

	destination := unmarshalFrame.Destination.String()

	return nil, destination
}

func Encode(fromTAP []byte) (error, string) {
	toTCP := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, toTCP)

	if _, err := encoder.Write(fromTAP); err != nil {
		return err, ""
	}

	if err := encoder.Close(); err != nil {
		return err, ""
	}

	return nil, toTCP.String() + "\n"
}

func Decode(fromTCP string) (error, []byte) {
	fromTCPWithoutNewline := strings.TrimSuffix(fromTCP, "\n")

	toTAP, err := base64.StdEncoding.DecodeString(fromTCPWithoutNewline)

	return err, toTAP
}
