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
	expectedReadChan := make(chan [wrappers.WrappedFrameSize]byte)
	expectedRemoteAddr, listener, err := getListenAddrWithFreePort()
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
				expectedReadChan,
				expectedRemoteAddr,
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

			timeout := time.Millisecond * 100
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
