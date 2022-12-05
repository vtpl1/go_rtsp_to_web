package bits

import "testing"

func TestGolombBitReader_ReadBit(t *testing.T) {
	tests := []struct {
		name    string
		self    *GolombBitReader
		wantRes uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.self.ReadBit()
			if (err != nil) != tt.wantErr {
				t.Errorf("GolombBitReader.ReadBit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("GolombBitReader.ReadBit() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGolombBitReader_ReadBits(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		self    *GolombBitReader
		args    args
		wantRes uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.self.ReadBits(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GolombBitReader.ReadBits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("GolombBitReader.ReadBits() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGolombBitReader_ReadBits32(t *testing.T) {
	type args struct {
		n uint
	}
	tests := []struct {
		name    string
		self    *GolombBitReader
		args    args
		wantR   uint32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := tt.self.ReadBits32(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GolombBitReader.ReadBits32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("GolombBitReader.ReadBits32() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestGolombBitReader_ReadBits64(t *testing.T) {
	type args struct {
		n uint
	}
	tests := []struct {
		name    string
		self    *GolombBitReader
		args    args
		wantR   uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := tt.self.ReadBits64(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GolombBitReader.ReadBits64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("GolombBitReader.ReadBits64() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func TestGolombBitReader_ReadExponentialGolombCode(t *testing.T) {
	tests := []struct {
		name    string
		self    *GolombBitReader
		wantRes uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.self.ReadExponentialGolombCode()
			if (err != nil) != tt.wantErr {
				t.Errorf("GolombBitReader.ReadExponentialGolombCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("GolombBitReader.ReadExponentialGolombCode() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGolombBitReader_ReadSE(t *testing.T) {
	tests := []struct {
		name    string
		self    *GolombBitReader
		wantRes uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.self.ReadSE()
			if (err != nil) != tt.wantErr {
				t.Errorf("GolombBitReader.ReadSE() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRes != tt.wantRes {
				t.Errorf("GolombBitReader.ReadSE() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
