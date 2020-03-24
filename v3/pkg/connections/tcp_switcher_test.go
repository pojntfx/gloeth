package connections

import (
	"log"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/pojntfx/gloeth/v3/pkg/wrappers"
)

func getListenAddrWithFreePort() (*net.TCPAddr, *net.TCPListener, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, err
	}

	return l.Addr().(*net.TCPAddr), l.(*net.TCPListener), err
}

func waitForConnection(listener *net.TCPListener) error {
	if _, err := listener.AcceptTCP(); err != nil {
		return err
	}

	return nil
}

func getConnection(raddr *net.TCPAddr) (*net.TCPConn, error) {
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func getFrame() [wrappers.WrappedFrameSize]byte {
	return [wrappers.WrappedFrameSize]byte{1}
}

func writeFrame(conn *net.TCPConn, frame [wrappers.WrappedFrameSize]byte) error {
	_, err := conn.Write(frame[:])

	return err
}

func readFrame(conn *net.TCPConn) ([wrappers.WrappedFrameSize]byte, error) {
	frame := [wrappers.WrappedFrameSize]byte{}

	_, err := conn.Read(frame[:])
	if err != nil {
		return [wrappers.WrappedFrameSize]byte{}, err
	}

	return frame, nil
}

func TestNewTCPSwitcher(t *testing.T) {
	expectedReadChan := make(chan [wrappers.WrappedFrameSize]byte)
	expectedRemoteAddr, _, err := getListenAddrWithFreePort()
	if err != nil {
		t.Error(err)
	}

	type args struct {
		readChan   chan [wrappers.WrappedFrameSize]byte
		remoteAddr *net.TCPAddr
	}
	tests := []struct {
		name string
		args args
		want *TCPSwitcher
	}{
		{
			"New",
			args{
				expectedReadChan,
				expectedRemoteAddr,
			},
			&TCPSwitcher{
				expectedReadChan,
				expectedRemoteAddr,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTCPSwitcher(tt.args.readChan, tt.args.remoteAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTCPSwitcher() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCPSwitcher_Open(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	remoteAddr, listener, err := getListenAddrWithFreePort()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan   chan [wrappers.WrappedFrameSize]byte
		remoteAddr *net.TCPAddr
		conn       *net.TCPConn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Open",
			fields{
				readChan,
				remoteAddr,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCPSwitcher{
				readChan:   tt.fields.readChan,
				remoteAddr: tt.fields.remoteAddr,
				conn:       tt.fields.conn,
			}
			if err := s.Open(); (err != nil) != tt.wantErr {
				t.Errorf("TCPSwitcher.Open() error = %v, wantErr %v", err, tt.wantErr)
			}

			timeout := time.Millisecond * 10
			timeoutChan := make(chan bool)
			go func() {
				time.Sleep(timeout)

				timeoutChan <- true
			}()

			go func() {
				if err := waitForConnection(listener); err != nil {
					t.Errorf("TCPSwitcher.Open() did not connect to TCP server: %v", err)
				}

				timeoutChan <- false
			}()

			if <-timeoutChan {
				log.Fatalf("TCPSwitcher.Open() did not connect to TCP server within %v", timeout)
			}
		})
	}
}

func TestTCPSwitcher_Close(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	remoteAddr, _, err := getListenAddrWithFreePort()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConnection(remoteAddr)

	type fields struct {
		readChan   chan [wrappers.WrappedFrameSize]byte
		remoteAddr *net.TCPAddr
		conn       *net.TCPConn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Close",
			fields{
				readChan:   readChan,
				conn:       conn,
				remoteAddr: remoteAddr,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCPSwitcher{
				readChan:   tt.fields.readChan,
				remoteAddr: tt.fields.remoteAddr,
				conn:       tt.fields.conn,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("TCPSwitcher.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTCPSwitcher_Read(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	remoteAddr, listener, err := getListenAddrWithFreePort()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConnection(remoteAddr)
	expectedFrame := getFrame()

	type fields struct {
		readChan   chan [wrappers.WrappedFrameSize]byte
		remoteAddr *net.TCPAddr
		conn       *net.TCPConn
	}
	tests := []struct {
		name               string
		fields             fields
		frameToWrite, want [wrappers.WrappedFrameSize]byte
		framesToTransceive uint
		wantErr            bool
	}{
		{
			"Read",
			fields{
				readChan:   readChan,
				conn:       conn,
				remoteAddr: remoteAddr,
			},
			expectedFrame,
			expectedFrame,
			5,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCPSwitcher{
				readChan:   tt.fields.readChan,
				remoteAddr: tt.fields.remoteAddr,
				conn:       tt.fields.conn,
			}

			go func() {
				if err := s.Read(); (err != nil) != tt.wantErr {
					t.Errorf("TCPSwitcher.Read() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			go func() {
				conn, err := listener.AcceptTCP()
				if err != nil {
					t.Errorf("TCPSwitcher.Read() TCP server mock error = %vv", err)
				}

				time.Sleep(time.Millisecond * 5)

				for i := 0; i < int(tt.framesToTransceive); i++ {
					if err := writeFrame(conn, tt.frameToWrite); err != nil {
						t.Error(err)
					}
				}
			}()

			for i := 0; i < int(tt.framesToTransceive); i++ {
				got := <-readChan

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("TCPSwitcher.Read() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestTCPSwitcher_Write(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	remoteAddr, listener, err := getListenAddrWithFreePort()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConnection(remoteAddr)
	expectedFrame := getFrame()

	type fields struct {
		readChan   chan [wrappers.WrappedFrameSize]byte
		remoteAddr *net.TCPAddr
		conn       *net.TCPConn
	}
	type args struct {
		frame [wrappers.WrappedFrameSize]byte
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		framesToTransceive uint
		want               [wrappers.WrappedFrameSize]byte
		wantErr            bool
	}{
		{
			"Write",
			fields{
				readChan,
				remoteAddr,
				conn,
			},
			args{
				expectedFrame,
			},
			5,
			expectedFrame,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCPSwitcher{
				readChan:   tt.fields.readChan,
				remoteAddr: tt.fields.remoteAddr,
				conn:       tt.fields.conn,
			}

			doneChan := make(chan bool)

			go func() {
				time.Sleep(time.Millisecond * 5)

				conn, err := listener.AcceptTCP()
				if err != nil {
					t.Error(err)
				}

				for i := 0; i < int(tt.framesToTransceive); i++ {
					got, err := readFrame(conn)
					if err != nil {
						t.Error(err)
					}

					if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("read(TCPSwitcher.Write()) = %v, want %v", got, tt.want)
					}
				}

				doneChan <- true
			}()

			for i := 0; i < int(tt.framesToTransceive); i++ {
				if err := s.Write(tt.args.frame); (err != nil) != tt.wantErr {
					t.Errorf("TCPSwitcher.Write() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			<-doneChan
		})
	}
}
