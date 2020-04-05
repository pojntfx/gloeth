package connections

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"testing"
	"time"

	gm "github.com/cseeger-epages/mac-gen-go"
	"github.com/pojntfx/gloeth/pkg/switchers"
)

func getMACAddress() (net.HardwareAddr, error) {
	prefix := gm.GenerateRandomLocalMacPrefix(false)
	suffix, err := gm.CalculateNICSufix(net.ParseIP("10.0.0.1"))
	if err != nil {
		return nil, err
	}

	rawDest := fmt.Sprintf("%v:%v", prefix, suffix)

	return net.ParseMAC(rawDest)
}

func TestNewSwitcherInfo(t *testing.T) {
	expectedReadChan := make(chan *net.HardwareAddr)
	expectedRaddr, _, err := getListener()
	if err != nil {
		t.Error(err)
	}

	type args struct {
		readChan chan *net.HardwareAddr
		raddr    *net.TCPAddr
	}
	tests := []struct {
		name string
		args args
		want *SwitcherInfo
	}{
		{
			"New",
			args{
				expectedReadChan,
				expectedRaddr,
			},
			&SwitcherInfo{
				expectedReadChan,
				expectedRaddr,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSwitcherInfo(tt.args.readChan, tt.args.raddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSwitcherInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwitcherInfo_Read(t *testing.T) {
	readChan := make(chan *net.HardwareAddr)
	raddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConn(raddr)
	if err != nil {
		t.Error(err)
	}
	expectedMAC, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan chan *net.HardwareAddr
		raddr    *net.TCPAddr
		conn     *net.TCPConn
	}
	tests := []struct {
		name             string
		fields           fields
		amountOfRequests uint
		macToWrite       *net.HardwareAddr
		want             *net.HardwareAddr
		wantErr          bool
	}{
		{
			"Read",
			fields{
				readChan,
				raddr,
				conn,
			},
			5,
			&expectedMAC,
			&expectedMAC,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SwitcherInfo{
				readChan: tt.fields.readChan,
				raddr:    tt.fields.raddr,
				conn:     tt.fields.conn,
			}
			go func() {
				if err := s.Read(); (err != nil) != tt.wantErr {
					t.Errorf("SwitcherInfo.Read() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			go func() {
				conn, err := listener.AcceptTCP()
				if err != nil {
					t.Error(err)
				}

				for i := 0; i < int(tt.amountOfRequests); i++ {
					outMAC := [switchers.SwitcherInfoSize]byte{}
					copy(outMAC[:], tt.macToWrite.String())

					if _, err := conn.Write(outMAC[:]); err != nil {
						t.Error(err)
					}
				}
			}()

			for i := 0; i < int(tt.amountOfRequests); i++ {
				got := <-readChan

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("SwitcherInfo.Read() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestSwitcherInfo_Open(t *testing.T) {
	raddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan chan *net.HardwareAddr
		raddr    *net.TCPAddr
		conn     *net.TCPConn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Open",
			fields{
				nil,
				raddr,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SwitcherInfo{
				readChan: tt.fields.readChan,
				raddr:    tt.fields.raddr,
				conn:     tt.fields.conn,
			}
			if err := s.Open(); (err != nil) != tt.wantErr {
				t.Errorf("SwitcherInfo.Open() error = %v, wantErr %v", err, tt.wantErr)
			}

			timeout := time.Millisecond * 10
			timeoutChan := make(chan bool)
			go func() {
				time.Sleep(timeout)

				timeoutChan <- true
			}()

			go func() {
				if err := waitForConn(listener); err != nil {
					t.Errorf("SwitcherInfo.Open() did not connect to TCP server: %v", err)
				}

				timeoutChan <- false
			}()

			if <-timeoutChan {
				log.Fatalf("SwitcherInfo.Open() did not connect to TCP server within %v", timeout)
			}
		})
	}
}

func TestSwitcherInfo_Close(t *testing.T) {
	raddr, _, err := getListener()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConn(raddr)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan chan *net.HardwareAddr
		raddr    *net.TCPAddr
		conn     *net.TCPConn
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Close",
			fields{
				nil,
				nil,
				conn,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SwitcherInfo{
				readChan: tt.fields.readChan,
				raddr:    tt.fields.raddr,
				conn:     tt.fields.conn,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("SwitcherInfo.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
