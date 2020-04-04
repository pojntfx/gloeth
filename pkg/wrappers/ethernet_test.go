package wrappers

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	gm "github.com/cseeger-epages/mac-gen-go"
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

func getHops() ([HopsCount]*net.HardwareAddr, error) {
	outHops := [HopsCount]*net.HardwareAddr{}

	for i := 0; i < HopsCount; i++ {
		mac, err := getMACAddress()
		if err != nil {
			return [HopsCount]*net.HardwareAddr{}, err
		}

		outHops[i] = &mac
	}

	return outHops, nil
}

func getFrame() [EncryptedFrameSize]byte {
	return [EncryptedFrameSize]byte{1}
}

func getWrappedFrame(dest, src *net.HardwareAddr, hops [HopsCount]*net.HardwareAddr, frame [EncryptedFrameSize]byte) [WrappedFrameSize]byte {
	outFrame := [WrappedFrameSize]byte{}

	outDest := [HeaderDestSize]byte{}
	copy(outDest[:], dest.String())

	outSrc := [HeaderSrcSize]byte{}
	copy(outSrc[:], src.String())

	outHops := [HopsSize]byte{}
	for i, hop := range hops {
		if hop == nil {
			continue
		}

		outHop := [HopSize]byte{}
		copy(outHop[:], hop.String())

		copy(outHops[i*HopSize:(i+1)*HopSize], outHop[:])
	}

	outHeader := [HeaderSize]byte{}
	copy(outHeader[:HeaderDestSize], outDest[:])
	copy(outHeader[HeaderDestSize:HeaderDestSize+HeaderSrcSize], outSrc[:])
	copy(outHeader[HeaderDestSize+HeaderSrcSize:HeaderDestSize+HeaderSrcSize+HopsSize], outHops[:])

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
	expectedDest, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	expectedSrc, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	expectedHops, err := getHops()
	if err != nil {
		t.Error(err)
	}
	expectedEmptyHops := [HopsCount]*net.HardwareAddr{}
	wrappedFrame := getWrappedFrame(&expectedDest, &expectedSrc, expectedHops, expectedFrame)
	wrappedFrameWithEmptyHops := getWrappedFrame(&expectedDest, &expectedSrc, expectedEmptyHops, expectedFrame)

	type args struct {
		dest  *net.HardwareAddr
		src   *net.HardwareAddr
		hops  [HopsCount]*net.HardwareAddr
		frame [EncryptedFrameSize]byte
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
				&expectedDest,
				&expectedSrc,
				expectedHops,
				expectedFrame,
			},
			wrappedFrame,
			false,
		},
		{
			"Wrap (empty hops)",
			NewEthernet(),
			args{
				&expectedDest,
				&expectedSrc,
				expectedEmptyHops,
				expectedFrame,
			},
			wrappedFrameWithEmptyHops,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ethernet{}
			got, err := e.Wrap(tt.args.dest, tt.args.src, tt.args.hops, tt.args.frame)
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
	expectedDest, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	expectedSrc, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	expectedHops, err := getHops()
	if err != nil {
		t.Error(err)
	}
	expectedEmptyHops := [HopsCount]*net.HardwareAddr{}
	wrappedFrame := getWrappedFrame(&expectedDest, &expectedSrc, expectedHops, expectedFrame)
	wrappedFrameWithEmptyHops := getWrappedFrame(&expectedDest, &expectedSrc, expectedEmptyHops, expectedFrame)

	type args struct {
		frame [WrappedFrameSize]byte
	}
	tests := []struct {
		name    string
		e       *Ethernet
		args    args
		want    *net.HardwareAddr
		want1   *net.HardwareAddr
		want2   [HopsCount]*net.HardwareAddr
		want3   [EncryptedFrameSize]byte
		wantErr bool
	}{
		{
			"Unwrap",
			NewEthernet(),
			args{
				wrappedFrame,
			},
			&expectedDest,
			&expectedSrc,
			expectedHops,
			expectedFrame,
			false,
		},
		{
			"Unwrap (empty hops)",
			NewEthernet(),
			args{
				wrappedFrameWithEmptyHops,
			},
			&expectedDest,
			&expectedSrc,
			expectedEmptyHops,
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
			nil,
			[HopsCount]*net.HardwareAddr{},
			[EncryptedFrameSize]byte{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ethernet{}
			got, got1, got2, got3, err := e.Unwrap(tt.args.frame)
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
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("Ethernet.Unwrap() got2 = %v, want %v", got2, tt.want2)
			}
			if !reflect.DeepEqual(got3, tt.want3) {
				t.Errorf("Ethernet.Unwrap() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func TestEthernet_GetShiftedHops(t *testing.T) {
	mac1, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	mac2, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}
	inHops := [HopsCount]*net.HardwareAddr{&mac1, &mac2}
	inHopsEmpty := [HopsCount]*net.HardwareAddr{}
	expectedHops := [HopsCount]*net.HardwareAddr{&mac2}

	type args struct {
		hops [HopsCount]*net.HardwareAddr
	}
	tests := []struct {
		name string
		e    *Ethernet
		args args
		want [HopsCount]*net.HardwareAddr
	}{
		{
			"GetShiftedHops",
			NewEthernet(),
			args{
				inHops,
			},
			expectedHops,
		},
		{
			"GetShiftedHops (empty hops)",
			NewEthernet(),
			args{
				inHopsEmpty,
			},
			inHopsEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Ethernet{}
			if got := e.GetShiftedHops(tt.args.hops); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ethernet.GetShiftedHops() = %v, want %v", got, tt.want)
			}
		})
	}
}
