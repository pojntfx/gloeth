package devices

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"testing"

	"github.com/pojntfx/gloeth/v3/pkg/encryptors"
	"github.com/vishvananda/netlink"
)

func getDevName(seed int) string {
	return fmt.Sprintf("testtap%v", seed+rand.Intn(99))
}

func TestNewTAP(t *testing.T) {
	readChan := make(chan [encryptors.PlaintextFrameSize]byte)
	mtu := uint(MTU)
	name := getDevName(100)

	type args struct {
		readChan chan [encryptors.PlaintextFrameSize]byte
		mtu      uint
		name     string
	}
	tests := []struct {
		name string
		args args
		want *TAP
	}{
		{
			"New",
			args{
				readChan,
				mtu,
				name,
			},
			&TAP{
				readChan,
				mtu,
				name,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTAP(tt.args.readChan, tt.args.mtu, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTAP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAP_Open(t *testing.T) {
	if os.Geteuid() != 0 && !testing.Short() {
		t.Skip()
	}

	readChan := make(chan [encryptors.PlaintextFrameSize]byte)
	mtu := uint(MTU)
	name := getDevName(200)

	type fields struct {
		readChan chan [encryptors.PlaintextFrameSize]byte
		mtu      uint
		name     string
	}
	tests := []struct {
		name         string
		fields       fields
		expectedName string
		expectedMTU  uint
		wantErr      bool
	}{
		{
			"Open",
			fields{
				readChan,
				mtu,
				name,
			},
			name,
			mtu,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TAP{
				readChan: tt.fields.readChan,
				mtu:      tt.fields.mtu,
				name:     tt.fields.name,
			}
			if err := s.Open(); (err != nil) != tt.wantErr {
				t.Errorf("TAP.Open() error = %v, wantErr %v", err, tt.wantErr)
			}

			link, err := netlink.LinkByName(tt.expectedName)
			if err != nil {
				t.Error(err)
			}
			if link == nil {
				t.Errorf("TAP.Open() link = %v, want !nil", link)
			}

			actualMTU := uint(link.Attrs().MTU)
			if actualMTU != tt.expectedMTU {
				t.Errorf("TAP.Open() mtu = %v, want %v", actualMTU, tt.expectedMTU)
			}

			if err := s.dev.Close(); err != nil {
				t.Error(err)
			}
		})
	}
}
