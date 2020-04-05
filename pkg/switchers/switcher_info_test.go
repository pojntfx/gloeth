package switchers

import (
	"net"
	"reflect"
	"testing"
)

func TestNewSwitcherInfo(t *testing.T) {
	tests := []struct {
		name string
		want *SwitcherInfo
	}{
		{
			"NewSwitcherInfo",
			&SwitcherInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSwitcherInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSwitcherInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwitcherInfo_RequestMACAddress(t *testing.T) {
	type fields struct {
		mac *net.HardwareAddr
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
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SwitcherInfo{
				mac: tt.fields.mac,
			}
			if err := s.RequestMACAddress(); (err != nil) != tt.wantErr {
				t.Errorf("SwitcherInfo.RequestMACAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := net.ParseMAC(s.mac.String()); err != nil {
				t.Errorf("parseMAC(SwitcherInfo.RequestMACAddress()) error = %v", err)
			}
		})
	}
}
