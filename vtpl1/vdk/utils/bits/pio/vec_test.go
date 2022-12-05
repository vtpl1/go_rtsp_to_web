package pio

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVecLen(t *testing.T) {
	type args struct {
		vec [][]byte
	}
	tests := []struct {
		name  string
		args  args
		wantN int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := VecLen(tt.args.vec); gotN != tt.wantN {
				t.Errorf("VecLen() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestVecSliceTo(t *testing.T) {
	type args struct {
		in  [][]byte
		out [][]byte
		s   int
		e   int
	}
	tests := []struct {
		name  string
		args  args
		wantN int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotN := VecSliceTo(tt.args.in, tt.args.out, tt.args.s, tt.args.e); gotN != tt.wantN {
				t.Errorf("VecSliceTo() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestVecSlice(t *testing.T) {
	type args struct {
		in [][]byte
		s  int
		e  int
	}
	tests := []struct {
		name    string
		args    args
		wantOut [][]byte
	}{
		{
			"1,-1",
			args{[][]byte{{1, 2, 3}, {4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}}, 1, -1},
			[][]byte{{2, 3}, {4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}},
		},
		{
			"2,-1",
			args{[][]byte{{1, 2, 3}, {4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}}, 2, -1},
			[][]byte{{3}, {4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}},
		},
		{
			"8,8",
			args{[][]byte{{1, 2, 3}, {4, 5, 6, 7, 8, 9}, {10, 11, 12, 13}}, 8, 8},
			[][]byte{{},{},{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := VecSlice(tt.args.in, tt.args.s, tt.args.e);
			fmt.Println(gotOut, tt.wantOut)
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("VecSlice() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
