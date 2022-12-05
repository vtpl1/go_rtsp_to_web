package bits

import "testing"

func TestReader_ReadBits64(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name     string
		self     *Reader
		args     args
		wantBits uint64
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBits, err := tt.self.ReadBits64(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadBits64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBits != tt.wantBits {
				t.Errorf("Reader.ReadBits64() = %v, want %v", gotBits, tt.wantBits)
			}
		})
	}
}

func TestReader_ReadBits(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name     string
		self     *Reader
		args     args
		wantBits uint
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBits, err := tt.self.ReadBits(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadBits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBits != tt.wantBits {
				t.Errorf("Reader.ReadBits() = %v, want %v", gotBits, tt.wantBits)
			}
		})
	}
}

func TestReader_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		self    *Reader
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.self.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Reader.Read() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestWriter_WriteBits64(t *testing.T) {
	type args struct {
		bits uint64
		n    int
	}
	tests := []struct {
		name    string
		self    *Writer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.self.WriteBits64(tt.args.bits, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteBits64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriter_WriteBits(t *testing.T) {
	type args struct {
		bits uint
		n    int
	}
	tests := []struct {
		name    string
		self    *Writer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.self.WriteBits(tt.args.bits, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteBits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		self    *Writer
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.self.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Writer.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Writer.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestWriter_FlushBits(t *testing.T) {
	tests := []struct {
		name    string
		self    *Writer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.self.FlushBits(); (err != nil) != tt.wantErr {
				t.Errorf("Writer.FlushBits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/*
func TestBits(t *testing.T) {
	rdata := []byte{0xf3, 0xb3, 0x45, 0x60}
	rbuf := bytes.NewReader(rdata[:])
	r := &Reader{R: rbuf}
	var u32 uint
	if u32, _ = r.ReadBits(4); u32 != 0xf {
		t.FailNow()
	}
	if u32, _ = r.ReadBits(4); u32 != 0x3 {
		t.FailNow()
	}
	if u32, _ = r.ReadBits(2); u32 != 0x2 {
		t.FailNow()
	}
	if u32, _ = r.ReadBits(2); u32 != 0x3 {
		t.FailNow()
	}
	b := make([]byte, 2)
	if r.Read(b); b[0] != 0x34 || b[1] != 0x56 {
		t.FailNow()
	}

	wbuf := &bytes.Buffer{}
	w := &Writer{W: wbuf}
	w.WriteBits(0xf, 4)
	w.WriteBits(0x3, 4)
	w.WriteBits(0x2, 2)
	w.WriteBits(0x3, 2)
	n, _ := w.Write([]byte{0x34, 0x56})
	if n != 2 {
		t.FailNow()
	}
	w.FlushBits()
	wdata := wbuf.Bytes()
	if wdata[0] != 0xf3 || wdata[1] != 0xb3 || wdata[2] != 0x45 || wdata[3] != 0x60 {
		t.FailNow()
	}

	// b = make([]byte, 8)
	// PutUInt64BE(b, 0x11223344, 32)
	// if b[0] != 0x11 || b[1] != 0x22 || b[2] != 0x33 || b[3] != 0x44 {
	// 	t.FailNow()
	// }
}

func TestReader_ReadBits64(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name     string
		self     *Reader
		args     args
		wantBits uint64
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBits, err := tt.self.ReadBits64(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadBits64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBits != tt.wantBits {
				t.Errorf("Reader.ReadBits64() = %v, want %v", gotBits, tt.wantBits)
			}
		})
	}
}

func TestReader_ReadBits(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name     string
		self     *Reader
		args     args
		wantBits uint
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBits, err := tt.self.ReadBits(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadBits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBits != tt.wantBits {
				t.Errorf("Reader.ReadBits() = %v, want %v", gotBits, tt.wantBits)
			}
		})
	}
}

func TestReader_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		self    *Reader
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.self.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Reader.Read() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestWriter_WriteBits64(t *testing.T) {
	type args struct {
		bits uint64
		n    int
	}
	tests := []struct {
		name    string
		self    *Writer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.self.WriteBits64(tt.args.bits, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteBits64() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriter_WriteBits(t *testing.T) {
	type args struct {
		bits uint
		n    int
	}
	tests := []struct {
		name    string
		self    *Writer
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.self.WriteBits(tt.args.bits, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteBits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriter_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		self    *Writer
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.self.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Writer.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Writer.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestWriter_FlushBits(t *testing.T) {
	tests := []struct {
		name    string
		self    *Writer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.self.FlushBits(); (err != nil) != tt.wantErr {
				t.Errorf("Writer.FlushBits() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

*/