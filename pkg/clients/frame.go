package clients

import (
	"context"

	proto "github.com/pojntfx/gloeth/pkg/proto/generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type FrameClient struct {
	remoteAddress     string
	remoteCertificate string
	channel           proto.FrameService_TransceiveFramesClient
}

func NewFrameClient(remoteAddress string, remoteCertificate string) *FrameClient {
	return &FrameClient{remoteAddress, remoteCertificate, nil}
}

func (c *FrameClient) Open() error {
	creds, err := credentials.NewClientTLSFromFile(c.remoteCertificate, "")
	if err != nil {
		return err
	}

	connection, err := grpc.Dial(c.remoteAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return err
	}

	client := proto.NewFrameServiceClient(connection)

	channel, err := client.TransceiveFrames(context.Background())
	if err != nil {
		return err
	}

	c.channel = channel

	return nil
}

func (c *FrameClient) Write(frame *proto.FrameMessage) error {
	return c.channel.Send(frame)
}

func (c *FrameClient) Read() (*proto.FrameMessage, error) {
	return c.channel.Recv()
}
