package routing

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	gm "github.com/cseeger-epages/mac-gen-go"
	"github.com/sauerbraten/graph/v2"
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

func getGraphFromMACs(mac1, mac2 *net.HardwareAddr, between []*net.HardwareAddr, alt []*net.HardwareAddr) *graph.Graph {
	rawGraph := [][2]string{}

	if len(between) == 0 {
		rawGraph = append(rawGraph, [2]string{mac1.String(), mac2.String()})
	}

	if len(between) > 1 {
		for i, hop := range between {
			if i == 0 {
				rawGraph = append(rawGraph, [2]string{mac1.String(), hop.String()})

				continue
			}

			rawGraph = append(rawGraph, [2]string{between[i-1].String(), hop.String()})

			if i == (len(between) - 1) {
				rawGraph = append(rawGraph, [2]string{hop.String(), mac2.String()})
			}
		}
	}

	if len(alt) > 1 {
		for i, hop := range alt {
			if i == 0 {
				rawGraph = append(rawGraph, [2]string{mac1.String(), hop.String()})

				continue
			}

			rawGraph = append(rawGraph, [2]string{alt[i-1].String(), hop.String()})
		}
	}

	return GetGraphFromRawData(rawGraph)
}

func TestNewRoutingTable(t *testing.T) {
	tests := []struct {
		name string
		want *RoutingTable
	}{
		{
			"New",
			&RoutingTable{
				graph.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRoutingTable(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoutingTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoutingTable_Register(t *testing.T) {
	mac1, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	mac2, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	inGraph := graph.New()
	expectedOutGraph := [][2]string{{mac1.String(), mac2.String()}}

	type fields struct {
		graph *graph.Graph
	}
	type args struct {
		mac1 *net.HardwareAddr
		mac2 *net.HardwareAddr
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    [][2]string
		wantErr bool
	}{
		{
			"Register",
			fields{
				inGraph,
			},
			args{
				&mac1,
				&mac2,
			},
			expectedOutGraph,
			false,
		},
		{
			"Register (same node)",
			fields{
				inGraph,
			},
			args{
				&mac1,
				&mac1,
			},
			expectedOutGraph,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				graph: tt.fields.graph,
			}
			if err := r.Register(tt.args.mac1, tt.args.mac2); (err != nil) != tt.wantErr {
				t.Errorf("RoutingTable.Register() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := GetRawDataFromGraph(r.graph)

			actualMatchLength := 0
			expectedMatchLength := len(tt.want)
			for _, ael := range got {
				for _, eel := range tt.want {
					if (ael[0] == eel[1] && ael[1] == eel[0]) || (ael[0] == eel[0] && ael[1] == eel[1]) {
						actualMatchLength = actualMatchLength + 1
					}
				}
			}

			if actualMatchLength != expectedMatchLength {
				t.Errorf("RoutingTable.Register() error = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoutingTable_GetHops(t *testing.T) {
	mac1, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	mac2, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	mac3, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	mac4, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	mac5, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	mac6, err := getMACAddress()
	if err != nil {
		t.Error(err)
	}

	expectedHops1 := []*net.HardwareAddr{&mac3, &mac4}
	inGraph1 := getGraphFromMACs(&mac1, &mac2, expectedHops1, []*net.HardwareAddr{&mac5, &mac6})

	expectedHops2 := []*net.HardwareAddr{}
	inGraph2 := getGraphFromMACs(&mac1, &mac2, expectedHops2, []*net.HardwareAddr{})

	inGraph3 := getGraphFromMACs(&mac1, &mac2, expectedHops1, []*net.HardwareAddr{})

	type fields struct {
		graph *graph.Graph
	}
	type args struct {
		switcherMAC *net.HardwareAddr
		adapterMAC  *net.HardwareAddr
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*net.HardwareAddr
		wantErr bool
	}{
		{
			"GetHops",
			fields{
				inGraph1,
			},
			args{
				&mac1,
				&mac2,
			},
			expectedHops1,
			false,
		},
		{
			"GetHops (direct connection)",
			fields{
				inGraph2,
			},
			args{
				&mac1,
				&mac2,
			},
			expectedHops2,
			false,
		},
		{
			"GetHops (no alt)",
			fields{
				inGraph3,
			},
			args{
				&mac1,
				&mac2,
			},
			expectedHops1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RoutingTable{
				graph: tt.fields.graph,
			}
			got, err := r.GetHops(tt.args.switcherMAC, tt.args.adapterMAC)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoutingTable.GetHops() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoutingTable.GetHops() = %v, want %v", got, tt.want)
			}
		})
	}
}
