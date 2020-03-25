package devices

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"reflect"
	"testing"

	"github.com/mdlayher/raw"
	"github.com/pojntfx/ethernet"
	"github.com/pojntfx/gloeth/pkg/encryptors"
	"github.com/songgao/water"
	"github.com/vishvananda/netlink"
)

func getDevName(seed int) string {
	return fmt.Sprintf("testtap%v", seed+rand.Intn(99))
}

func getDev(name string, mtu uint) (*water.Interface, error) {
	dev, err := water.New(water.Config{
		DeviceType: water.TAP,
		PlatformSpecificParams: water.PlatformSpecificParams{
			Name: name,
		},
	})
	if err != nil {
		return nil, err
	}

	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil, err
	}

	if err := netlink.LinkSetMTU(link, int(mtu)); err != nil {
		return nil, err
	}

	if err := netlink.LinkSetUp(link); err != nil {
		return nil, err
	}

	return dev, nil
}

func closeDev(dev *TAP) error {
	dev.closed = true

	return dev.dev.Close()
}

func writeTestFrame(devName, content string) error {
	const etherType = 0xcccc
	dest := ethernet.Broadcast

	link, err := net.InterfaceByName(devName)
	if err != nil {
		return err
	}

	conn, err := raw.ListenPacket(link, etherType, nil)
	if err != nil {
		return err
	}

	frame := &ethernet.Frame{
		Destination: dest,
		Source:      link.HardwareAddr,
		EtherType:   etherType,
		Payload:     []byte(content),
	}

	ethFrame, err := frame.MarshalBinary()
	if err != nil {
		return err
	}

	outFrame := [encryptors.PlaintextFrameSize]byte{}
	copy(outFrame[:], ethFrame)

	_, err = conn.WriteTo(outFrame[:], &raw.Addr{
		HardwareAddr: dest,
	})

	return err
}

func readFrame(frame [encryptors.PlaintextFrameSize]byte) (ethernet.Frame, error) {
	var ethFrame ethernet.Frame
	err := ethFrame.UnmarshalBinary(frame[:])

	return ethFrame, err
}

func getFrame() [encryptors.PlaintextFrameSize]byte {
	return [encryptors.PlaintextFrameSize]byte{1}
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
				true,
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
		closed   bool
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
				true,
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
				closed:   tt.fields.closed,
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

			if err := closeDev(s); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestTAP_Close(t *testing.T) {
	if os.Geteuid() != 0 && !testing.Short() {
		t.Skip()
	}

	readChan := make(chan [encryptors.PlaintextFrameSize]byte)
	mtu := uint(MTU)
	name := getDevName(200)
	dev, err := getDev(name, mtu)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		readChan chan [encryptors.PlaintextFrameSize]byte
		mtu      uint
		name     string
		dev      *water.Interface
		closed   bool
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
				mtu,
				name,
				dev,
				false,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TAP{
				readChan: tt.fields.readChan,
				mtu:      tt.fields.mtu,
				name:     tt.fields.name,
				dev:      tt.fields.dev,
				closed:   tt.fields.closed,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("TAP.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s.closed != true {
				t.Errorf("TAP.Close() closed = %v, wantErr %v", s.closed, true)
			}
		})
	}
}

func TestTAP_Read(t *testing.T) {
	if os.Geteuid() != 0 && !testing.Short() {
		t.Skip()
	}

	readChan := make(chan [encryptors.PlaintextFrameSize]byte)
	mtu := uint(MTU)
	name := getDevName(300)
	dev, err := getDev(name, mtu)
	if err != nil {
		t.Error(err)
	}
	expectedContent := "test"

	type fields struct {
		readChan chan [encryptors.PlaintextFrameSize]byte
		mtu      uint
		name     string
		dev      *water.Interface
		closed   bool
	}
	tests := []struct {
		name                 string
		fields               fields
		contentToWrite, want string
		framesToTransceive   uint
		wantErr              bool
	}{
		{
			"Read",
			fields{
				readChan,
				mtu,
				name,
				dev,
				false,
			},
			expectedContent,
			expectedContent,
			5,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TAP{
				readChan: tt.fields.readChan,
				mtu:      tt.fields.mtu,
				name:     tt.fields.name,
				dev:      tt.fields.dev,
				closed:   tt.fields.closed,
			}

			go func() {
				if err := s.Read(); (err != nil) != tt.wantErr {
					t.Errorf("TAP.Read() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			doneChan := make(chan bool)
			go func() {
				for i := 0; i < int(tt.framesToTransceive); i++ {
					if err := writeTestFrame(s.name, tt.contentToWrite); err != nil {
						t.Error(err)
					}
				}

				doneChan <- true
			}()

			for matches := 0; matches < int(tt.framesToTransceive); matches++ {
				frame := <-readChan

				inFrame, err := readFrame(frame)
				if err != nil {
					t.Error(frame)
				}

				actualContent := string(inFrame.Payload[:len(tt.contentToWrite)])

				if actualContent == tt.want {
					matches = matches + 1

					continue
				}

				matches = matches - 1
			}

			<-doneChan
			if err := closeDev(s); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestTAP_Write(t *testing.T) {
	if os.Geteuid() != 0 && !testing.Short() {
		t.Skip()
	}

	readChan := make(chan [encryptors.PlaintextFrameSize]byte)
	mtu := uint(MTU)
	name := getDevName(400)
	dev, err := getDev(name, mtu)
	if err != nil {
		t.Error(err)
	}
	frameToWrite := getFrame()

	type fields struct {
		readChan chan [encryptors.PlaintextFrameSize]byte
		mtu      uint
		name     string
		dev      *water.Interface
		closed   bool
	}
	type args struct {
		frame [encryptors.PlaintextFrameSize]byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Write",
			fields{
				readChan,
				mtu,
				name,
				dev,
				false,
			},
			args{
				frameToWrite,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &TAP{
				readChan: tt.fields.readChan,
				mtu:      tt.fields.mtu,
				name:     tt.fields.name,
				dev:      tt.fields.dev,
				closed:   tt.fields.closed,
			}

			if err := s.Write(tt.args.frame); (err != nil) != tt.wantErr {
				t.Errorf("TAP.Write() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := closeDev(s); err != nil {
				t.Error(err)
			}
		})
	}
}
