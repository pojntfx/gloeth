package transceiver

import (
	"bufio"
	"github.com/pojntfx/gloeth/pkg/protocol"
	"log"
	"net"
)

type TCP struct {
	SendHostPort   string
	ListenHostPort string
}

func (t *TCP) Send(frame protocol.Frame) error {
	log.Println("tcp transceiver sending frame", frame)

	return nil
}

func (t *TCP) Listen(errors chan error, receivedFrames chan protocol.Frame) {
	log.Println("tcp transceiver listening")

	listener, err := net.Listen("tcp", t.ListenHostPort)
	if err != nil {
		errors <- err
	}
	defer func() {
		if err := listener.Close(); err != nil {
			errors <- err
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			errors <- err
		}

		go t.handleConnection(conn, errors, receivedFrames)
	}
}

func (t *TCP) handleConnection(conn net.Conn, errors chan error, receivedFrames chan protocol.Frame) {
	message, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		if err := conn.Close(); err != nil {
			errors <- err
		}

		return
	}

	var frame protocol.Frame
	if err := frame.Unserialize(message); err != nil {
		errors <- err
	}

	receivedFrames <- frame

	if _, err := conn.Write([]byte("\n")); err != nil {
		errors <- err
	}

	t.handleConnection(conn, errors, receivedFrames)
}
