package routing

import (
	"reflect"
	"testing"
)

func TestGetDifferenceOfNestedArrays(t *testing.T) {
	old := [][2]string{
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

	new := [][2]string{
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

	expectedDeletions := [][2]string{
		{
			"b",
			"c",
		},
		{
			"c",
			"d",
		},
	}

	expectedAdditions := [][2]string{
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
		old [][2]string
		new [][2]string
	}
	tests := []struct {
		name          string
		args          args
		wantDeletions [][2]string
		wantAdditions [][2]string
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
	in1 := [][2]string{
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

	in2 := [][2]string{
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

	expectedOut := []string{
		"b",
		"a",
		"c",
	}

	type args struct {
		in [][2]string
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
			expectedOut,
		},
		{
			"GetUniqueKeys (different order)",
			args{
				in2,
			},
			expectedOut,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUniqueKeys(tt.args.in)

			actualLen := len(got)
			expectedLen := len(tt.want)

			actualMatchLength := 0
			expectedMatchLength := len(tt.want)
			for _, ael := range got {
				for _, eel := range tt.want {
					if ael == eel {
						actualMatchLength = actualMatchLength + 1
					}
				}
			}

			if actualLen != expectedLen {
				t.Errorf("len(GetUniqueKeys()) = %v, want %v", actualLen, expectedLen)
			}

			if actualMatchLength != expectedMatchLength {
				t.Errorf("len(matches(GetUniqueKeys())) = %v, want %v", actualMatchLength, expectedMatchLength)
			}
		})
	}
}

func TestDeduplicateNestedArray(t *testing.T) {
	in1 := [][2]string{
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

	in2 := [][2]string{
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

	expectedOut := [][2]string{
		{
			"b", // Note the potentially reversed value; this order is not guaranteed
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
	}

	type args struct {
		in [][2]string
	}
	tests := []struct {
		name string
		args args
		want [][2]string
	}{
		{
			"DeduplicateNestedArray",
			args{
				in1,
			},
			expectedOut,
		},
		{
			"DeduplicateNestedArray (different order)",
			args{
				in2,
			},
			expectedOut,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeduplicateNestedArray(tt.args.in)

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
				t.Errorf("len(DeduplicateNestedArray()) = %v, want %v", actualLen, expectedLen)
			}

			if actualMatchLength != expectedMatchLength {
				t.Errorf("len(matches(DeduplicateNestedArray())) = %v, want %v", actualMatchLength, expectedMatchLength)
			}
		})
	}
}
