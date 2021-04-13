package services

import (
	proto "github.com/pojntfx/gloeth/pkg/proto/generated"
)

//go:generate sh -c "mkdir -p ../proto/generated && protoc --go_out=paths=source_relative,plugins=grpc:../proto/generated -I=../proto ../proto/*.proto"

type FrameService struct {
	proto.UnimplementedFrameServiceServer
	channel proto.FrameService_TransceiveFramesServer
}

func NewFrameService() *FrameService {
	return &FrameService{}
}

func (s *FrameService) TransceiveFrames(channel proto.FrameService_TransceiveFramesServer) error {
	wait := make(chan struct{})

	s.channel = channel

	<-wait // Wait forever

	return nil
}

func (s *FrameService) Write(frame *proto.FrameMessage) error {
	s.waitTillOpen()

	return s.channel.Send(frame)
}

func (s *FrameService) Read() (*proto.FrameMessage, error) {
	s.waitTillOpen()

	return s.channel.Recv()
}

func (s *FrameService) waitTillOpen() {
	for s.channel == nil {

	}
}
