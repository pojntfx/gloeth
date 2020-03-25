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

func getFrame(destMAC, srcMAC net.HardwareAddr) ([encryptors.PlaintextFrameSize]byte, error) {
	var frame ethernet.Frame
	frame.Destination = destMAC
	frame.Source = srcMAC

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
	destMAC, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	srcMAC, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	frame, err := getFrame(destMAC, srcMAC)
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
		want1   *net.HardwareAddr
		want2   *net.HardwareAddr
		wantErr bool
	}{
		{
			"GetDestMACAddress",
			NewEthernet(),
			args{
				frame,
			},
			&destMAC,
			&srcMAC,
			false,
		},
		{
			"GetDestMACAddress (faulty frame)",
			NewEthernet(),
			args{
				faultyFrame,
			},
			&expectedMacFromFaultyFrame,
			&expectedMacFromFaultyFrame,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ethernet{}
			got1, got2, err := e.GetMACAddresses(tt.args.frame)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ethernet.GetDestMACAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Ethernet.GetDestMACAddress()[0] = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("Ethernet.GetDestMACAddress()[1] = %v, want %v", got2, tt.want2)
			}
		})
	}
}
