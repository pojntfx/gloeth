package encoders

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	gm "github.com/cseeger-epages/mac-gen-go"
	"github.com/pojntfx/ethernet"
	"github.com/pojntfx/gloeth/v3/pkg/encryptors"
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

func getFrame(mac net.HardwareAddr) ([encryptors.PlaintextFrameSize]byte, error) {
	var frame ethernet.Frame
	frame.Destination = mac

	rawOutFrame, err := frame.MarshalBinary()
	if err != nil {
		return [encryptors.PlaintextFrameSize]byte{}, nil
	}

	outFrame := [encryptors.PlaintextFrameSize]byte{}
	copy(outFrame[:], rawOutFrame)

	return outFrame, nil
}

func getFaultyFrame() [encryptors.PlaintextFrameSize]byte {
	return [encryptors.PlaintextFrameSize]byte{}
}

func TestNewEthernet(t *testing.T) {
	tests := []struct {
		name string
		want *Ethernet
	}{
		{
			"New",
			&Ethernet{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEthernet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEthernet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEthernet_GetDestMACAddress(t *testing.T) {
	mac, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	frame, err := getFrame(mac)
	if err != nil {
		t.Error(err)
	}
	faultyFrame := getFaultyFrame()
	expectedMacFromFaultyFrame, err := net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		t.Error(err)
	}

	type args struct {
		frame [encryptors.PlaintextFrameSize]byte
	}
	tests := []struct {
		name    string
		e       *Ethernet
		args    args
		want    *net.HardwareAddr
		wantErr bool
	}{
		{
			"GetDestMACAddress",
			NewEthernet(),
			args{
				frame,
			},
			&mac,
			false,
		},
		{
			"GetDestMACAddress (faulty frame)",
			NewEthernet(),
			args{
				faultyFrame,
			},
			&expectedMacFromFaultyFrame,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ethernet{}
			got, err := e.GetDestMACAddress(tt.args.frame)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ethernet.GetDestMACAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ethernet.GetDestMACAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
