package pkg

import (
	"errors"
	"strings"
)

const (
	FromStart = "f="
	ToStart   = "t="
	BodyStart = "b="

	Delimiter = "|"
)

type Frame struct {
	From, To string
	Body     []byte
}

func Unserialize(payload string) (error, Frame) {
	if !strings.Contains(payload, FromStart) {
		return errors.New("missing from"), Frame{}
	}

	if !strings.Contains(payload, ToStart) {
		return errors.New("missing to"), Frame{}
	}

	if !strings.Contains(payload, BodyStart) {
		return errors.New("missing body"), Frame{}
	}

	if !(len(strings.Split(payload, Delimiter)) == 4) {
		return errors.New("missing delimiter"), Frame{}
	}

	from := strings.Split(strings.Split(payload, FromStart)[1], Delimiter)[0]
	to := strings.Split(strings.Split(payload, ToStart)[1], Delimiter)[0]
	body := []byte(strings.Split(strings.Split(payload, BodyStart)[1], Delimiter)[0])

	return nil, Frame{
		From: from,
		To:   to,
		Body: body,
	}
}

func Serialize(frame Frame) (error, string) {
	return nil, FromStart + frame.From + Delimiter + ToStart + frame.To + Delimiter + BodyStart + string(frame.Body) + Delimiter
}
