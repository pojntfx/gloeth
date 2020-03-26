package routing

import (
	"reflect"
	"testing"
)

func TestGetDifferenceOfNestedArrays(t *testing.T) {
	old := [][]string{
		{
			"a",
			"b",
		},
		{
			"b",
			"c",
		},
		{
			"c",
			"d",
		},
	}

	new := [][]string{
		{
			"a",
			"b",
		},
		{
			"b",
			"d",
		},
		{
			"d",
			"e",
		},
	}

	expectedDeletions := [][]string{
		{
			"b",
			"c",
		},
		{
			"c",
			"d",
		},
	}

	expectedAdditions := [][]string{
		{
			"b",
			"d",
		},
		{
			"d",
			"e",
		},
	}

	type args struct {
		old [][]string
		new [][]string
	}
	tests := []struct {
		name          string
		args          args
		wantDeletions [][]string
		wantAdditions [][]string
	}{
		{
			"GetDifferenceOfNestedArrays",
			args{
				old,
				new,
			},
			expectedDeletions,
			expectedAdditions,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDeletions, gotAdditions := GetDifferenceOfNestedArrays(tt.args.old, tt.args.new)
			if !reflect.DeepEqual(gotDeletions, tt.wantDeletions) {
				t.Errorf("GetDifferenceOfNestedArrays() gotDeletions = %v, want %v", gotDeletions, tt.wantDeletions)
			}
			if !reflect.DeepEqual(gotAdditions, tt.wantAdditions) {
				t.Errorf("GetDifferenceOfNestedArrays() gotAdditions = %v, want %v", gotAdditions, tt.wantAdditions)
			}
		})
	}
}

func TestGetUniqueKeys(t *testing.T) {
	in1 := [][]string{
		{
			"b",
			"a",
		},
		{
			"b",
			"c",
		},
		{
			"a",
			"c",
		},
		{
			"c",
			"a",
		},
		{
			"a",
			"b",
		},
		{
			"b",
			"a",
		},
	}

	in2 := [][]string{
		{
			"a",
			"b",
		},
		{
			"b",
			"c",
		},
		{
			"a",
			"c",
		},
		{
			"c",
			"a",
		},
		{
			"b",
			"a",
		},
		{
			"b",
			"a",
		},
	}

	expectedOut1 := []string{
		"b",
		"a",
		"c",
	}
	expectedOut2 := []string{
		"a",
		"b",
		"c",
	}

	type args struct {
		in [][]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"GetUniqueKeys",
			args{
				in1,
			},
			expectedOut1,
		},
		{
			"GetUniqueKeys (different order)",
			args{
				in2,
			},
			expectedOut2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUniqueKeys(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUniqueKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
