package servers

import (
	"net"

	proto "github.com/pojntfx/gloeth/pkg/proto/generated"
	"github.com/pojntfx/gloeth/pkg/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type FrameServer struct {
	listenAddress string
	certificate   string
	key           string
	frameService  *services.FrameService
}

func NewFrameServer(listenAddress string, certificate string, key string, frameService *services.FrameService) *FrameServer {
	return &FrameServer{listenAddress, certificate, key, frameService}
}

func (s *FrameServer) Open() error {
	listenAddress, err := net.ResolveTCPAddr("tcp", s.listenAddress)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", listenAddress)
	if err != nil {
		return err
	}

	creds, err := credentials.NewServerTLSFromFile(s.certificate, s.key)
	if err != nil {
		return err
	}

	server := grpc.NewServer(grpc.Creds(creds))

	reflection.Register(server)
	proto.RegisterFrameServiceServer(server, s.frameService)

	if err := server.Serve(listener); err != nil {
		return err
	}

	return nil
}
