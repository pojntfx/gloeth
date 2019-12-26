package server

import (
	"bufio"
	"fmt"
	"github.com/pojntfx/gloeth/pkg/protocol"
	"net"
)

const (
	OkResponse               = "0\n"
	UnserializeErrorResponse = "1\n"
	UnknownErrorResponse     = "1000\n"
)

func HandleConnection(connection net.Conn) {
	data, err := bufio.NewReader(connection).ReadBytes('\n')
	if err != nil {
		connection.Close()
		return
	}

	var frame protocol.Frame

	if err := frame.Unserialize(string(data)); err != nil {
		if err == err.(*protocol.UnserializeError) {
			connection.Write([]byte(UnserializeErrorResponse))
			return
		}

		connection.Write([]byte(UnknownErrorResponse))
		return
	}

	fmt.Println(frame)

	connection.Write([]byte(OkResponse))

	HandleConnection(connection)
}
