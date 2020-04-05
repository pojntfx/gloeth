package switchers

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func readSwitcherInfo(conn *net.TCPConn) (net.HardwareAddr, error) {
	rawMAC := [SwitcherInfoSize]byte{}

	if _, err := conn.Read(rawMAC[:]); err != nil {
		return nil, err
	}

	return net.ParseMAC(string(rawMAC[:]))
}

func TestNewSwitcherInfo(t *testing.T) {
	expectedLaddr, _, err := getListener()
	if err != nil {
		t.Error(err)
	}

	type args struct {
		laddr *net.TCPAddr
	}
	tests := []struct {
		name string
		args args
		want *SwitcherInfo
	}{
		{
			"New",
			args{
				expectedLaddr,
			},
			&SwitcherInfo{
				expectedLaddr,
				nil,
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSwitcherInfo(tt.args.laddr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSwitcherInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwitcherInfo_RequestMACAddress(t *testing.T) {
	type fields struct {
		laddr    *net.TCPAddr
		listener *net.TCPListener
		mac      *net.HardwareAddr
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"RequestMACAddress",
			fields{
				nil,
				nil,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SwitcherInfo{
				laddr: tt.fields.laddr,
				mac:   tt.fields.mac,
			}
			if err := s.RequestMACAddress(); (err != nil) != tt.wantErr {
				t.Errorf("SwitcherInfo.RequestMACAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err := s.RequestMACAddress(); (err != nil) != tt.wantErr {
				t.Errorf("SwitcherInfo.RequestMACAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSwitcherInfo_Open(t *testing.T) {
	laddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}
	if err := listener.Close(); err != nil {
		t.Error(err)
	}
	mac, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		laddr    *net.TCPAddr
		listener *net.TCPListener
		mac      *net.HardwareAddr
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Open",
			fields{
				laddr,
				nil,
				&mac,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SwitcherInfo{
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				mac:      tt.fields.mac,
			}
			if err := s.Open(); (err != nil) != tt.wantErr {
				t.Errorf("SwitcherInfo.Open() error = %v, wantErr %v", err, tt.wantErr)
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
				t.Errorf("SwitcherInfo.Open() did not connect to TCP client within %v", timeout)
			}
		})
	}
}

func TestSwitcherInfo_Close(t *testing.T) {
	_, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		laddr    *net.TCPAddr
		listener *net.TCPListener
		mac      *net.HardwareAddr
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
				listener,
				nil,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SwitcherInfo{
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				mac:      tt.fields.mac,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("SwitcherInfo.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSwitcherInfo_Read(t *testing.T) {
	laddr, listener, err := getListener()
	if err != nil {
		t.Error(err)
	}
	conn, err := getConn(laddr)
	if err != nil {
		t.Error(err)
	}
	expectedMAC, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		laddr    *net.TCPAddr
		listener *net.TCPListener
		mac      *net.HardwareAddr
	}
	tests := []struct {
		name             string
		fields           fields
		amountOfRequests uint
		want             *net.HardwareAddr
		wantErr          bool
	}{
		{
			"Read",
			fields{
				nil,
				listener,
				&expectedMAC,
			},
			5,
			&expectedMAC,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SwitcherInfo{
				laddr:    tt.fields.laddr,
				listener: tt.fields.listener,
				mac:      tt.fields.mac,
			}
			go func() {
				if err := s.Read(); (err != nil) != tt.wantErr {
					t.Errorf("SwitcherInfo.Read() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()

			for i := 0; i < int(tt.amountOfRequests); i++ {
				got, err := readSwitcherInfo(conn)
				if err != nil {
					t.Errorf("parseMAC(read((SwitcherInfo.Read())) error = %v", err)
				}

				if !reflect.DeepEqual(&got, tt.want) {
					t.Errorf("SwitcherInfo.Read() = %v, want %v", got, tt.want)
				}

				if err := conn.Close(); err != nil {
					t.Error(err)
				}

				// The server does not close the connection; it could send a new MAC address in the future.
				// This is not necessary in the test, so close it.
				conn, err = getConn(laddr)
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}
