package connections

import (
	"fmt"
	"net"
	"reflect"
	"testing"

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
	expectedMAC, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan chan *net.HardwareAddr
		raddr    *net.TCPAddr
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
			}
			go func() {
				if err := s.Read(); (err != nil) != tt.wantErr {
					t.Errorf("SwitcherInfo.Read() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			go func() {
				for i := 0; i < int(tt.amountOfRequests); i++ {
					conn, err := listener.AcceptTCP()
					if err != nil {
						t.Error(err)
					}

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
