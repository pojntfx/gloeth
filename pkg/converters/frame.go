package converters

import proto "github.com/pojntfx/gloeth/pkg/proto/generated"

type FrameConverter struct {
}

func NewFrameConverter() *FrameConverter {
	return &FrameConverter{}
}

func (c *FrameConverter) ToExternal(rawFrame []byte) (proto.FrameMessage, error) {
	return proto.FrameMessage{Content: rawFrame}, nil
}

func (c *FrameConverter) ToInternal(frame proto.FrameMessage) ([]byte, error) {
	return frame.Content, nil
}
