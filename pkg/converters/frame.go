package converters

import proto "github.com/pojntfx/gloeth/pkg/proto/generated"

type FrameConverter struct {
}

func NewFrameConverter() *FrameConverter {
	return &FrameConverter{}
}

func (c *FrameConverter) ToExternal(rawFrame []byte, preSharedKey string) (*proto.FrameMessage, error) {
	return &proto.FrameMessage{Content: rawFrame, PreSharedKey: preSharedKey}, nil
}

func (c *FrameConverter) ToInternal(frame *proto.FrameMessage) ([]byte, string, error) {
	return frame.Content, frame.PreSharedKey, nil
}
