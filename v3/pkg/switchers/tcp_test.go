package switchers

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/pojntfx/gloeth/v3/pkg/wrappers"
)

func getListener() (*net.TCPAddr, *net.TCPListener, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, err
	}

	return l.Addr().(*net.TCPAddr), l.(*net.TCPListener), err
}

func getConn(raddr *net.TCPAddr) (*net.TCPConn, error) {
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

func TestNewTCP(t *testing.T) {
	expectedReadChan := make(chan [wrappers.WrappedFrameSize]byte)
	expectedLaddr, _, err := getListener()
	if err != nil {
		t.Error(err)
	}

	type args struct {
		readChan   chan [wrappers.WrappedFrameSize]byte
		listenAddr *net.TCPAddr
	}
	tests := []struct {
		name string
		args args
		want *TCP
	}{
		{
			"New",
			args{
				expectedReadChan,
				expectedLaddr,
			},
			&TCP{
				expectedReadChan,
				expectedLaddr,
				nil,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTCP(tt.args.readChan, tt.args.listenAddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTCP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTCP_Open(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	laddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}
	if err := listener.Close(); err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan   chan [wrappers.WrappedFrameSize]byte
		listenAddr *net.TCPAddr
		listener   *net.TCPListener
		conns      []map[string]*net.TCPConn
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
				laddr,
				nil,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				laddr:    tt.fields.listenAddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}
			if err := s.Open(); (err != nil) != tt.wantErr {
				t.Errorf("TCP.Open() error = %v, wantErr %v", err, tt.wantErr)
			}

			timeoutChan := make(chan bool)
			timeout := time.Millisecond * 10
			go func() {
				_, err := getConn(laddr)
				if err != nil {
					t.Error(err)
				}

				timeoutChan <- false
			}()

			go func() {
				time.Sleep(timeout)

				timeoutChan <- true
			}()

			if <-timeoutChan {
				t.Errorf("TCP.Open() did not connect to TCP client within %v", timeout)
			}
		})
	}
}

func TestTCP_Close(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	laddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan chan [wrappers.WrappedFrameSize]byte
		laddr    *net.TCPAddr
		listener *net.TCPListener
		conns    []map[string]*net.TCPConn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Close",
			fields{
				readChan,
				laddr,
				listener,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("TCP.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTCP_Read(t *testing.T) {
	readChan := make(chan [wrappers.WrappedFrameSize]byte)
	laddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConn(laddr)
	expectedFrame := getFrame()

	type fields struct {
		readChan chan [wrappers.WrappedFrameSize]byte
		laddr    *net.TCPAddr
		listener *net.TCPListener
		conns    []map[string]*net.TCPConn
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
				readChan,
				laddr,
				listener,
				nil,
			},
			expectedFrame,
			expectedFrame,
			5,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TCP{
				readChan: tt.fields.readChan,
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				conns:    tt.fields.conns,
			}

			go func() {
				if err := s.Read(); (err != nil) != tt.wantErr {
					t.Errorf("TCP.Read() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			go func() {
				time.Sleep(time.Millisecond * 5)

				for i := 0; i < int(tt.framesToTransceive); i++ {
					if err := writeFrame(conn, expectedFrame); err != nil {
						t.Error(err)
					}
				}
			}()

			for i := 0; i < int(tt.framesToTransceive); i++ {
				got := <-readChan

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("TCP.Read() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
