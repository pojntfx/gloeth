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

	deletions := [][]string{
		{
			"b",
			"c",
		},
		{
			"c",
			"d",
		},
	}

	additions := [][]string{
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
			deletions,
			additions,
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
