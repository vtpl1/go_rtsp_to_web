package pio

import "testing"

func TestU8(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := U8(tt.args.b); gotI != tt.wantI {
				t.Errorf("U8() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestU16BE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI uint16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := U16BE(tt.args.b); gotI != tt.wantI {
				t.Errorf("U16BE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestI16BE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI int16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := I16BE(tt.args.b); gotI != tt.wantI {
				t.Errorf("I16BE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestI24BE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := I24BE(tt.args.b); gotI != tt.wantI {
				t.Errorf("I24BE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestU24BE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := U24BE(tt.args.b); gotI != tt.wantI {
				t.Errorf("U24BE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestI32BE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI int32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := I32BE(tt.args.b); gotI != tt.wantI {
				t.Errorf("I32BE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestU32LE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := U32LE(tt.args.b); gotI != tt.wantI {
				t.Errorf("U32LE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestU32BE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := U32BE(tt.args.b); gotI != tt.wantI {
				t.Errorf("U32BE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestU40BE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := U40BE(tt.args.b); gotI != tt.wantI {
				t.Errorf("U40BE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestU64BE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := U64BE(tt.args.b); gotI != tt.wantI {
				t.Errorf("U64BE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}

func TestI64BE(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name  string
		args  args
		wantI int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotI := I64BE(tt.args.b); gotI != tt.wantI {
				t.Errorf("I64BE() = %v, want %v", gotI, tt.wantI)
			}
		})
	}
}
