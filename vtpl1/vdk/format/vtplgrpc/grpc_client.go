package vtplgrpc

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"html"
	"io"
	"math"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vtpl1/vdk/av"
	"github.com/vtpl1/vdk/codec/h264parser"
	"github.com/vtpl1/vdk/codec/h265parser"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var log = logrus.New()

func init() {
	// TODO: next add write to file
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	// log.SetLevel(Storage.ServerLogLevel())
}

const (
	SignalStreamGrpcStop = iota
	SignalGrpcCodecUpdate
)

type GRPCLient struct {
	headers map[string]string
	Signals chan int
	// OutgoingProxyQueue  chan *[]byte
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
	pURL                *url.URL
	username            string
	password            string
	control             string
	videoCodec          av.CodecType
	PreVideoTS          int64
	fuStarted           bool
	vps                 []byte
	sps                 []byte
	pps                 []byte
	// grpc_client         StreamServiceClient
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

func (client *GRPCLient) parseURL(rawURL string) error {
	l, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	username := l.User.Username()
	password, _ := l.User.Password()
	l.User = nil
	if l.Port() == "" {
		l.Host = fmt.Sprintf("%s:%s", l.Host, "554")
	}
	if l.Scheme != "grpc" {
		l.Scheme = "grpc"
	}
	client.pURL = l
	client.username = username
	client.password = password
	client.control = l.String()
	return nil
}

func Dial(options GRPCLientOptions) (*GRPCLient, error) {
	client := &GRPCLient{
		headers: make(map[string]string),
		Signals: make(chan int, 100),
		// OutgoingProxyQueue:  make(chan *[]byte, 3000),
		OutgoingPacketQueue: make(chan *av.Packet, 3000),
		BufferRtpPacket:     bytes.NewBuffer([]byte{}),
		videoID:             -1,
		audioID:             -2,
		videoIDX:            -1,
		audioIDX:            -2,
		options:             options,
		AudioTimeScale:      8000,
	}
	err := client.parseURL(html.UnescapeString(client.options.URL))
	if err != nil {
		return nil, err
	}
	go client.startStream()
	return client, nil
}

func (client *GRPCLient) Close() {
}

// binSize
func binSize(val int) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(val))
	return buf
}

//Println mini logging functions
func (client *GRPCLient) Println(v ...interface{}) {
	if client.options.Debug {
		fmt.Println(v)
	}
}

func (client *GRPCLient) CodecUpdateSPS(val []byte) {
	if client.videoCodec != av.H264 && client.videoCodec != av.H265 {
		return
	}
	if bytes.Compare(val, client.sps) == 0 {
		return
	}
	client.sps = val
	if (client.videoCodec == av.H264 && len(client.pps) == 0) || (client.videoCodec == av.H265 && (len(client.vps) == 0 || len(client.pps) == 0)) {
		return
	}
	var codecData av.VideoCodecData
	var err error
	switch client.videoCodec {
	case av.H264:
		client.Println("Codec Update SPS", val)
		codecData, err = h264parser.NewCodecDataFromSPSAndPPS(val, client.pps)
		if err != nil {
			client.Println("Parse Codec Data Error", err)
			return
		}
	case av.H265:
		codecData, err = h265parser.NewCodecDataFromVPSAndSPSAndPPS(client.vps, val, client.pps)
		if err != nil {
			client.Println("Parse Codec Data Error", err)
			return
		}
	}
	if len(client.CodecData) > 0 {
		for i, i2 := range client.CodecData {
			if i2.Type().IsVideo() {
				client.CodecData[i] = codecData
			}
		}
	} else {
		client.CodecData = append(client.CodecData, codecData)
	}
	client.Signals <- SignalGrpcCodecUpdate
}

func (client *GRPCLient) CodecUpdatePPS(val []byte) {
	if client.videoCodec != av.H264 && client.videoCodec != av.H265 {
		return
	}
	if bytes.Compare(val, client.pps) == 0 {
		return
	}
	client.pps = val
	if (client.videoCodec == av.H264 && len(client.sps) == 0) || (client.videoCodec == av.H265 && (len(client.vps) == 0 || len(client.sps) == 0)) {
		return
	}
	var codecData av.VideoCodecData
	var err error
	switch client.videoCodec {
	case av.H264:
		client.Println("Codec Update PPS", val)
		codecData, err = h264parser.NewCodecDataFromSPSAndPPS(client.sps, val)
		if err != nil {
			client.Println("Parse Codec Data Error", err)
			return
		}
	case av.H265:
		codecData, err = h265parser.NewCodecDataFromVPSAndSPSAndPPS(client.vps, client.sps, val)
		if err != nil {
			client.Println("Parse Codec Data Error", err)
			return
		}
	}
	if len(client.CodecData) > 0 {
		for i, i2 := range client.CodecData {
			if i2.Type().IsVideo() {
				client.CodecData[i] = codecData
			}
		}
	} else {
		client.CodecData = append(client.CodecData, codecData)
	}
	client.Signals <- SignalGrpcCodecUpdate
}

func (client *GRPCLient) startStream() {
	defer func() {
		client.Signals <- SignalStreamGrpcStop
	}()
	log.WithFields(logrus.Fields{
		"module": "vtplgrpc",
		"func":   "Dial",
	}).Infof("Dial: %v", client.pURL.Host)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(client.pURL.Host, opts...)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "vtplgrpc",
			"func":   "startStream",
		}).Fatalf("fail to Dial: %v", err)
	}
	defer conn.Close()
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "vtplgrpc",
			"func":   "Dial",
		}).Fatalf("fail to dial: %v", err)
	}
	grpc_client := NewStreamServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	
	defer cancel()

	read_frame_request := &ReadFrameRequest{}
	read_frame_request.ChannelId = 0
	read_frame_request.AppId = 0
	read_frame_request.MajorMinor = 0

	stream, err := grpc_client.ReadFrame(ctx, read_frame_request)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "vtplgrpc",
			"func":   "startStream",
		}).Fatalf("fail to startStream: %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.WithFields(logrus.Fields{
					"module": "vtplgrpc",
					"func":   "startStream",
				}).Fatalf("ReadFrame failed: %v", err)
			}
			var retmap []*av.Packet
			client.videoIDX = 1

			media_type := response.Frame.MediaType
			frame_type := response.Frame.FrameType
			time_stamp := (response.Frame.TimeStamp * 90) / 1000000
			frame_size := response.Frame.BufferSize
			bit_rate := response.Frame.BitRate
			fps := response.Frame.Fps
			motion_available := response.Frame.MotionAvailable
			stream_type := response.Frame.MajorMinor
			ssrc := response.Frame.Ssrc
			frame_id := response.Frame.FrameId
			content := response.Frame.Buffer
			offset := 0
			end := len(content)

			_ = frame_type
			_ = frame_size
			_ = bit_rate
			_ = fps
			_ = motion_available
			_ = stream_type
			_ = ssrc
			_ = frame_id


			if media_type == 2 {
				client.videoCodec = av.H264
				if client.PreVideoTS == 0 {
					client.PreVideoTS = time_stamp
				}
				if time_stamp-client.PreVideoTS < 0 {
					if math.MaxUint32-client.PreVideoTS < 90*100 { //100 ms
						client.PreVideoTS = 0
						client.PreVideoTS -= (math.MaxUint32 - client.PreVideoTS)
					} else {
						client.PreVideoTS = 0
					}
				}
				nalRaw := content
				nal := nalRaw[4:]

				naluType := nal[0] & 0x1f
				switch {
				case naluType >= 1 && naluType <= 5:
					retmap = append(retmap, &av.Packet{
						Data:            append(binSize(len(nal)), nal...),
						CompositionTime: time.Duration(1) * time.Millisecond,
						Idx:             client.videoIDX,
						IsKeyFrame:      naluType == 5,
						Duration:        time.Duration(float32(time_stamp-client.PreVideoTS)/90) * time.Millisecond,
						Time:            time.Duration(time_stamp/90) * time.Millisecond,
					})
				case naluType == 7:
					client.CodecUpdateSPS(nal)
				case naluType == 8:
					client.CodecUpdatePPS(nal)
					client.WaitCodec = true
				case naluType == 24:
					packet := nal[1:]
					for len(packet) >= 2 {
						size := int(packet[0])<<8 | int(packet[1])
						if size+2 > len(packet) {
							break
						}
						naluTypefs := packet[2] & 0x1f
						switch {
						case naluTypefs >= 1 && naluTypefs <= 5:
							retmap = append(retmap, &av.Packet{
								Data:            append(binSize(len(packet[2:size+2])), packet[2:size+2]...),
								CompositionTime: time.Duration(1) * time.Millisecond,
								Idx:             client.videoIDX,
								IsKeyFrame:      naluType == 5,
								Duration:        time.Duration(float32(time_stamp-client.PreVideoTS)/90) * time.Millisecond,
								Time:            time.Duration(time_stamp/90) * time.Millisecond,
							})
						case naluTypefs == 7:
							client.CodecUpdateSPS(packet[2 : size+2])
						case naluTypefs == 8:
							client.CodecUpdatePPS(packet[2 : size+2])
						}
						packet = packet[size+2:]
					}
				case naluType == 28:
					fuIndicator := content[offset]
					fuHeader := content[offset+1]
					isStart := fuHeader&0x80 != 0
					isEnd := fuHeader&0x40 != 0
					if isStart {
						client.fuStarted = true
						client.BufferRtpPacket.Truncate(0)
						client.BufferRtpPacket.Reset()
						client.BufferRtpPacket.Write([]byte{fuIndicator&0xe0 | fuHeader&0x1f})
					}
					if client.fuStarted {
						client.BufferRtpPacket.Write(content[offset+2 : end])
						if isEnd {
							client.fuStarted = false
							naluTypef := client.BufferRtpPacket.Bytes()[0] & 0x1f
							if naluTypef == 7 || naluTypef == 9 {
								bufered, _ := h264parser.SplitNALUs(append([]byte{0, 0, 0, 1}, client.BufferRtpPacket.Bytes()...))
								for _, v := range bufered {
									naluTypefs := v[0] & 0x1f
									switch {
									case naluTypefs == 5:
										client.BufferRtpPacket.Reset()
										client.BufferRtpPacket.Write(v)
										naluTypef = 5
									case naluTypefs == 7:
										client.CodecUpdateSPS(v)
									case naluTypefs == 8:
										client.CodecUpdatePPS(v)
									}
								}
							}
							retmap = append(retmap, &av.Packet{
								Data:            append(binSize(client.BufferRtpPacket.Len()), client.BufferRtpPacket.Bytes()...),
								CompositionTime: time.Duration(1) * time.Millisecond,
								Duration:        time.Duration(float32(time_stamp-client.PreVideoTS)/90) * time.Millisecond,
								Idx:             client.videoIDX,
								IsKeyFrame:      naluTypef == 5,
								Time:            time.Duration(time_stamp/90) * time.Millisecond,
							})
						}
					}
				default:
					log.WithFields(logrus.Fields{
						"module": "vtplgrpc",
						"func":   "startStream",
					}).Infof("Unsupported NAL Type %d", naluType)
				}
			}
			if len(retmap) > 0 {
				client.PreVideoTS = time_stamp		
				for _, i2 := range retmap {
					if len(client.OutgoingPacketQueue) > 2000 {
						client.Println("RTSP Client OutgoingPacket Chanel Full")
						return
					}
					client.OutgoingPacketQueue <- i2
				}		
			}
		}
	}()
	stream.CloseSend()
	<-waitc
}
