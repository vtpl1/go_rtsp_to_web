package pio

import "testing"

func TestPutU8(t *testing.T) {
	type args struct {
		b []byte
		v uint8
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutU8(tt.args.b, tt.args.v)
		})
	}
}

func TestPutI16BE(t *testing.T) {
	type args struct {
		b []byte
		v int16
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutI16BE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutU16BE(t *testing.T) {
	type args struct {
		b []byte
		v uint16
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutU16BE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutI24BE(t *testing.T) {
	type args struct {
		b []byte
		v int32
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutI24BE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutU24BE(t *testing.T) {
	type args struct {
		b []byte
		v uint32
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutU24BE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutI32BE(t *testing.T) {
	type args struct {
		b []byte
		v int32
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutI32BE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutU32BE(t *testing.T) {
	type args struct {
		b []byte
		v uint32
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutU32BE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutU32LE(t *testing.T) {
	type args struct {
		b []byte
		v uint32
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutU32LE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutU40BE(t *testing.T) {
	type args struct {
		b []byte
		v uint64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutU40BE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutU48BE(t *testing.T) {
	type args struct {
		b []byte
		v uint64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutU48BE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutU64BE(t *testing.T) {
	type args struct {
		b []byte
		v uint64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutU64BE(tt.args.b, tt.args.v)
		})
	}
}

func TestPutI64BE(t *testing.T) {
	type args struct {
		b []byte
		v int64
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PutI64BE(tt.args.b, tt.args.v)
		})
	}
}
