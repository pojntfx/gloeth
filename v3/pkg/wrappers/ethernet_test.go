package wrappers

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	gm "github.com/cseeger-epages/mac-gen-go"
)

func getDest() (net.HardwareAddr, error) {
	prefix := gm.GenerateRandomLocalMacPrefix(false)
	suffix, err := gm.CalculateNICSufix(net.ParseIP("10.0.0.1"))
	if err != nil {
		return nil, err
	}

	rawDest := fmt.Sprintf("%v:%v", prefix, suffix)

	return net.ParseMAC(rawDest)
}

func getFrame() [FrameSize]byte {
	return [FrameSize]byte{1}
}

func getWrappedFrame(dest net.HardwareAddr, frame [FrameSize]byte) [WrappedFrameSize]byte {
	outFrame := [WrappedFrameSize]byte{}

	outDest := [DestSize]byte{}
	copy(outDest[:], dest.String())

	outHeader := [HeaderSize]byte{}
	copy(outHeader[:DestSize], outDest[:])

	copy(outFrame[:HeaderSize], outHeader[:])
	copy(outFrame[HeaderSize:], frame[:])

	return outFrame
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

func TestEthernet_Wrap(t *testing.T) {
	expectedFrame := getFrame()
	expectedDest, err := getDest()
	if err != nil {
		t.Error(err)
	}
	wrappedFrame := getWrappedFrame(expectedDest, expectedFrame)

	type args struct {
		dest  *net.HardwareAddr
		frame [FrameSize]byte
	}
	tests := []struct {
		name    string
		e       *Ethernet
		args    args
		want    [WrappedFrameSize]byte
		wantErr bool
	}{
		{
			"Wrap",
			NewEthernet(),
			args{
				dest:  &expectedDest,
				frame: expectedFrame,
			},
			wrappedFrame,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ethernet{}
			got, err := e.Wrap(tt.args.dest, tt.args.frame)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ethernet.Wrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ethernet.Wrap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEthernet_Unwrap(t *testing.T) {
	expectedFrame := getFrame()
	expectedDest, err := getDest()
	if err != nil {
		t.Error(err)
	}
	wrappedFrame := getWrappedFrame(expectedDest, expectedFrame)

	type args struct {
		frame [WrappedFrameSize]byte
	}
	tests := []struct {
		name    string
		e       *Ethernet
		args    args
		want    *net.HardwareAddr
		want1   [FrameSize]byte
		wantErr bool
	}{
		{
			"Unwrap",
			NewEthernet(),
			args{
				wrappedFrame,
			},
			&expectedDest,
			expectedFrame,
			false,
		},
		{
			"Unwrap (faulty frame)",
			NewEthernet(),
			args{
				[WrappedFrameSize]byte{},
			},
			nil,
			[FrameSize]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ethernet{}
			got, got1, err := e.Unwrap(tt.args.frame)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ethernet.Unwrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ethernet.Unwrap() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Ethernet.Unwrap() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
