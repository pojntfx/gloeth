package services

import proto "github.com/pojntfx/gloeth/pkg/proto/generated"

//go:generate sh -c "mkdir -p ../proto/generated && protoc --go_out=paths=source_relative,plugins=grpc:../proto/generated -I=../proto ../proto/*.proto"

type FrameService struct {
	proto.UnimplementedFrameServiceServer
}

func NewFrameService() *FrameService {
	return &FrameService{}
}
