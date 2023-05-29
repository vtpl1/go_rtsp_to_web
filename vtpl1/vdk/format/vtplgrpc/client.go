package vtplgrpc

import (
	"bytes"
	"time"

	"github.com/vtpl1/vdk/av"
)

type GRPCLient struct {
	headers             map[string]string
	Signals             chan int
	OutgoingProxyQueue  chan *[]byte
	OutgoingPacketQueue chan *av.Packet
	BufferRtpPacket     *bytes.Buffer
	videoID             int
	audioID             int
	videoIDX            int8
	audioIDX            int8
	options             GRPCLientOptions
	AudioTimeScale      int64
	WaitCodec           bool
	CodecData           []av.CodecData
	SDPRaw              []byte
}

type GRPCLientOptions struct {
	Debug              bool
	URL                string
	DialTimeout        time.Duration
	ReadWriteTimeout   time.Duration
	DisableAudio       bool
	OutgoingProxy      bool
	InsecureSkipVerify bool
}

func Dial(options GRPCLientOptions) (*GRPCLient, error) {
	client := &GRPCLient{
		headers:             make(map[string]string),
		Signals:             make(chan int, 100),
		OutgoingProxyQueue:  make(chan *[]byte, 3000),
		OutgoingPacketQueue: make(chan *av.Packet, 3000),
		BufferRtpPacket:     bytes.NewBuffer([]byte{}),
		videoID:             -1,
		audioID:             -2,
		videoIDX:            -1,
		audioIDX:            -2,
		options:             options,
		AudioTimeScale:      8000,
	}
	return client, nil
}

func (client *GRPCLient) Close() {
}
