package pkg

import (
	"bufio"
	"fmt"
	"net"
)

type TCP struct {
	ReadHostPort string
}

func (t *TCP) Write(errors chan error, status chan string, frame []byte, writeHostPort string) {
	status <- "writing frame to TCP transport"

	conn, err := net.Dial("tcp", writeHostPort)
	if err != nil {
		errors <- err
		t.Write(errors, status, frame, writeHostPort)
		return
	}

	err, encodedFrame := Encode(frame)
	if err != nil {
		errors <- err
		return
	}

	if _, err := fmt.Fprintf(conn, encodedFrame); err != nil {
		errors <- err
		t.Write(errors, status, frame, writeHostPort)
		return
	}

	if _, err := bufio.NewReader(conn).ReadBytes('\n'); err != nil {
		errors <- err
		t.Write(errors, status, frame, writeHostPort)
		return
	}

	if err := conn.Close(); err != nil {
		errors <- err
		t.Write(errors, status, frame, writeHostPort)
		return
	}

	status <- "wrote frame to TCP transport"
}

func (t *TCP) Read(errors chan error, status chan string, readFrames chan []byte) {
	status <- "reading frames from TCP transport"

	srv, err := net.Listen("tcp", t.ReadHostPort)
	if err != nil {
		errors <- err
		return
	}
	defer func(errors chan error, status chan string) {
		status <- "closing TCP transport"

		if err := srv.Close(); err != nil {
			errors <- err
			return
		}

		status <- "closed TCP transport"
	}(errors, status)

	for {
		status <- "reading frame from TCP transport"

		conn, err := srv.Accept()
		if err != nil {
			errors <- err
		}

		frame, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			errors <- err
		}

		if _, err := conn.Write([]byte("\n")); err != nil {
			errors <- err
		}

		err, decodedFrame := Decode(frame)
		if err != nil {
			errors <- err
		}

		readFrames <- decodedFrame

		status <- "read frame from TCP transport"
	}

	status <- "read frames from TCP transport"
}
