package routing

import (
	"testing"

	"github.com/sauerbraten/graph/v2"
)

func getRawGraphData() [][2]string {
	return [][2]string{
		{
			"n1",
			"n4",
		},
		{
			"n2",
			"n5",
		},
		{
			"n3",
			"n6",
		},
		{
			"n6",
			"n5",
		},
		{
			"n5",
			"n4",
		},
		{
			"n8",
			"n7",
		},
		{
			"n7",
			"n4",
		},
		{
			"n7",
			"n5",
		},
		{
			"n10",
			"n9",
		},
		{
			"n9",
			"n6",
		},
		{
			"n9",
			"n7",
		},
	}
}

func getGraph(rawData [][2]string) *graph.Graph {
	return GetGraphFromRawData(rawData)
}

func TestGetGraphFromRawData(t *testing.T) {
	in := getRawGraphData()

	type args struct {
		in [][2]string
	}
	tests := []struct {
		name string
		args args
		want [][2]string
	}{
		{
			"GetGraphFromRawData",
			args{
				in,
			},
			in,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRaw := GetGraphFromRawData(tt.args.in)

			got := GetRawDataFromGraph(gotRaw)

			actualLen := len(got)
			expectedLen := len(tt.want)

			actualMatchLength := 0
			expectedMatchLength := len(tt.want)
			for _, ael := range got {
				for _, eel := range tt.want {
					if (ael[0] == eel[1] && ael[1] == eel[0]) || (ael[0] == eel[0] && ael[1] == eel[1]) {
						actualMatchLength = actualMatchLength + 1
					}
				}
			}

			if actualLen != expectedLen {
				t.Errorf("len(GetGraphFromRawData()) = %v, want %v", actualLen, expectedLen)
			}

			if actualMatchLength != expectedMatchLength {
				t.Errorf("len(matches(GetGraphFromRawData())) = %v, want %v", actualMatchLength, expectedMatchLength)
			}
		})
	}
}
