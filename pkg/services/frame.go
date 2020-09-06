package services

import "github.com/pojntfx/gloeth/pkg/proto/generated/proto"

//go:generate sh -c "mkdir -p ../proto/generated && protoc --go_out=paths=source_relative,plugins=grpc:../proto/generated -I=../ ../proto/frame.proto"

type FrameService struct {
	proto.UnimplementedFrameServiceServer
}

func (s *FrameService) TransceiveFrames(server proto.FrameService_TransceiveFramesServer) error {
	frame, err := server.Recv()
	if err != nil {
		return err
	}

	if err := server.Send(frame); err != nil {
		return err
	}

	return nil
}
