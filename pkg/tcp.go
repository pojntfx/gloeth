package pkg

import (
	"bufio"
	"fmt"
	"net"
)

type TCP struct {
	WriteHostPort string
	ReadHostPort  string
}

func (t *TCP) Write(errors chan error, status chan string, frame []byte) {
	status <- "writing frame to TCP transport"

	conn, err := net.Dial("tcp", t.WriteHostPort)
	if err != nil {
		errors <- err
		t.Write(errors, status, frame)
		return
	}

	if _, err := fmt.Fprintf(conn, string(frame)+"\n"); err != nil {
		errors <- err
		t.Write(errors, status, frame)
		return
	}

	if _, err := bufio.NewReader(conn).ReadBytes('\n'); err != nil {
		errors <- err
		t.Write(errors, status, frame)
		return
	}

	if err := conn.Close(); err != nil {
		errors <- err
		t.Write(errors, status, frame)
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

		frame, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			errors <- err
		}

		if _, err := conn.Write([]byte("\n")); err != nil {
			errors <- err
		}

		readFrames <- frame

		status <- "read frame from TCP transport"
	}

	status <- "read frames from TCP transport"
}
