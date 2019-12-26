package protocol

import (
	"strings"
)

const (
	FromStart = "f="
	ToStart   = "t="
	BodyStart = "b="

	Delimiter = "|"
)

const (
	FrameMissingFromErrorMessage = "missing from"
	FrameMissingToErrorMessage   = "missing to"
	FrameMissingBodyErrorMessage = "missing body"

	FrameMissingDelimiterErrorMessage = "missing delimiter"
)

type Frame struct {
	From, To string
	Body     []byte
}

type UnserializeError struct {
	err   string
	field string
}

func (e *UnserializeError) Error() string {
	return e.err + " " + e.field
}

func (frame *Frame) Unserialize(payload string) error {
	if !strings.Contains(payload, FromStart) {
		return &UnserializeError{
			err:   FrameMissingFromErrorMessage,
			field: FromStart,
		}
	}

	if !strings.Contains(payload, ToStart) {
		return &UnserializeError{
			err:   FrameMissingToErrorMessage,
			field: ToStart,
		}
	}

	if !strings.Contains(payload, BodyStart) {
		return &UnserializeError{
			err:   FrameMissingBodyErrorMessage,
			field: BodyStart,
		}
	}

	if !(len(strings.Split(payload, Delimiter)) == 4) {
		return &UnserializeError{
			err:   FrameMissingDelimiterErrorMessage,
			field: Delimiter,
		}
	}

	from := strings.Split(strings.Split(payload, FromStart)[1], Delimiter)[0]
	to := strings.Split(strings.Split(payload, ToStart)[1], Delimiter)[0]
	body := []byte(strings.Split(strings.Split(payload, BodyStart)[1], Delimiter)[0])

	frame.From = from
	frame.To = to
	frame.Body = body

	return nil
}

func (frame *Frame) Serialize() (error, string) {
	return nil, FromStart + frame.From + Delimiter + ToStart + frame.To + Delimiter + BodyStart + string(frame.Body) + Delimiter
}
