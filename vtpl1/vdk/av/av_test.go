// Package av defines basic interfaces and data structures of container demux/mux and audio encode/decode.
package av

import (
	"testing"
)

func TestSampleFormat_BytesPerSample(t *testing.T) {
	tests := []struct {
		name string
		self SampleFormat
		want int
	}{
		{"U8", SampleFormat(U8), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.self.BytesPerSample(); got != tt.want {
				t.Errorf("SampleFormat.BytesPerSample() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestChannelLayout_String(t *testing.T) {
	tests := []struct {
		name string
		self ChannelLayout
		want string
	}{
		{"CH_MONO", ChannelLayout(CH_MONO), "1ch"},
		{"CH_STEREO", ChannelLayout(CH_STEREO), "2ch"},
		{"CH_2_1", ChannelLayout(CH_2_1), "3ch"},
		{"CH_2POINT1", ChannelLayout(CH_2POINT1), "3ch"},
		{"CH_SURROUND", ChannelLayout(CH_SURROUND), "3ch"},
		{"CH_3POINT1", ChannelLayout(CH_3POINT1), "4ch"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.self.String(); got != tt.want {
				t.Errorf("ChannelLayout.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
