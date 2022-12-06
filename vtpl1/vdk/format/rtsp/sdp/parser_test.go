package sdp

import (
	"reflect"
	"testing"

	"github.com/vtpl1/vdk/av"
)

func TestParse(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name       string
		args       args
		wantSess   Session
		wantMedias []Media
	}{
		{"parse-sdp-null", args{}, Session{}, nil},
		{"parse-sdp-1", args{
			`
v=0
o=- 1459325504777324 1 IN IP4 192.168.0.123
s=RTSP/RTP stream from Network Video Server
i=mpeg4cif
t=0 0
a=tool:LIVE555 Streaming Media v2009.09.28
a=type:broadcast
a=control:*
a=range:npt=0-
a=x-qt-text-nam:RTSP/RTP stream from Network Video Server
a=x-qt-text-inf:mpeg4cif
m=video 0 RTP/AVP 96
c=IN IP4 0.0.0.0
b=AS:300
a=rtpmap:96 H264/90000
a=fmtp:96 profile-level-id=420029; packetization-mode=1; sprop-parameter-sets=Z00AHpWoKA9k,aO48gA==
a=x-dimensions: 720, 480
a=x-framerate: 15
a=control:track1
m=audio 0 RTP/AVP 96
c=IN IP4 0.0.0.0
b=AS:256
a=rtpmap:96 MPEG4-GENERIC/16000/2
a=fmtp:96 streamtype=5;profile-level-id=1;mode=AAC-hbr;sizelength=13;indexlength=3;indexdeltalength=3;config=1408
a=control:track2
m=audio 0 RTP/AVP 0
c=IN IP4 0.0.0.0
b=AS:50
a=recvonly
a=control:rtsp://109.195.127.207:554/mpeg4cif/trackID=2
a=rtpmap:0 PCMU/8000
a=Media_header:MEDIAINFO=494D4B48010100000400010010710110401F000000FA000000000000000000000000000000000000;
a=appversion:1.0`,
		}, Session{}, []Media{
			{
				AVType:             "video",
				Type:               av.CodecType(av.H264),
				FPS:                15,
				TimeScale:          90000,
				Control:            "track1",
				Rtpmap:             96,
				ChannelCount:       0,
				Config:             nil,
				SpropParameterSets: [][]byte{{103, 77, 0, 30, 149, 168, 40, 15, 100}, {104, 238, 60, 128}},
				SpropVPS:           nil,
				SpropSPS:           nil,
				SpropPPS:           nil,
				PayloadType:        96,
				SizeLength:         0,
				IndexLength:        0,
			},
			{
				AVType:             "audio",
				Type:               av.CodecType(av.AAC),
				FPS:                0,
				TimeScale:          16000,
				Control:            "track2",
				Rtpmap:             96,
				ChannelCount:       0,
				Config:             []byte{20, 8},
				SpropParameterSets: nil,
				SpropVPS:           nil,
				SpropSPS:           nil,
				SpropPPS:           nil,
				PayloadType:        96,
				SizeLength:         13,
				IndexLength:        3,
			},
			{
				AVType:             "audio",
				Type:               av.CodecType(av.PCM_MULAW),
				FPS:                0,
				TimeScale:          8000,
				Control:            "rtsp://109.195.127.207:554/mpeg4cif/trackID=2",
				Rtpmap:             0,
				ChannelCount:       0,
				Config:             nil,
				SpropParameterSets: nil,
				SpropVPS:           nil,
				SpropSPS:           nil,
				SpropPPS:           nil,
				PayloadType:        0,
				SizeLength:         0,
				IndexLength:        0,
			},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSess, gotMedias := Parse(tt.args.content)
			if !reflect.DeepEqual(gotSess, tt.wantSess) {
				t.Errorf("Parse() gotSess = %v, want %v", gotSess, tt.wantSess)
			}
			if !reflect.DeepEqual(gotMedias, tt.wantMedias) {
				t.Errorf("Parse() gotMedias = %v, want %v", gotMedias, tt.wantMedias)
			}
		})
	}
}
