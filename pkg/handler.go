package pkg

import (
	"bufio"
	"fmt"
	"net"
)

func HandleConnection(connection net.Conn) {
	data, err := bufio.NewReader(connection).ReadBytes('\n')
	if err != nil {
		connection.Close()
		return
	}

	err, frame := Unserialize(string(data))
	if err != nil {
		connection.Write([]byte("1\n"))
		return
	}

	fmt.Println(frame)

	connection.Write([]byte("0\n"))

	HandleConnection(connection)
}
