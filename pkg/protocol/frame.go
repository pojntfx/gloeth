package protocol

import "strings"

type Frame struct {
	From string
	To   string
	Body []byte
}

const (
	FromStart = "f="
	ToStart   = "t="
	BodyStart = "b="

	Delimiter = "|"
)

func (frame *Frame) Serialize() (error, []byte) {
	return nil, []byte(FromStart + frame.From + Delimiter + ToStart + frame.To + Delimiter + BodyStart + string(frame.Body) + Delimiter)
}

func (frame *Frame) Unserialize(rawPayload []byte) error {
	payload := string(rawPayload)

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

type UnserializeError struct {
	err   string
	field string
}

const (
	FrameMissingFromErrorMessage = "missing from"
	FrameMissingToErrorMessage   = "missing to"
	FrameMissingBodyErrorMessage = "missing body"

	FrameMissingDelimiterErrorMessage = "missing delimiter"
)

func (e *UnserializeError) Error() string {
	return e.err + ": " + e.field
}
