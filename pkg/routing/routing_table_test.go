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
