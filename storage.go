package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/imdario/mergo"
	"github.com/vtpl1/go_rtsp_to_web/utils"
	"github.com/vtpl1/vdk/av"
)

type ClientST struct {
	mode              int
	signals           chan int
	outgoingAVPacket  chan *av.Packet
	outgoingRTPPacket chan *[]byte
	socket            net.Conn
}

type ChannelST struct {
	Name               string `json:"name,omitempty" groups:"api,config"`
	URL                string `json:"url,omitempty" groups:"api,config"`
	OnDemand           bool   `json:"on_demand,omitempty" groups:"api,config"`
	Debug              bool   `json:"debug,omitempty" groups:"api,config"`
	Status             int    `json:"status,omitempty" groups:"api"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify,omitempty" groups:"api,config"`
	Audio              bool   `json:"audio,omitempty" groups:"api,config"`
	runLock            bool
	codecs             []av.CodecData
	sdp                []byte
	signals            chan int
	// hlsSegmentBuffer   map[int]SegmentOld
	// hlsSegmentNumber   int
	clients map[string]ClientST
	ack     time.Time
	// hlsMuxer           *MuxerHLS `json:"-"`
}

// ServerST stream storage section
type StreamST struct {
	Name     string               `json:"name,omitempty" groups:"api,config"`
	Channels map[string]ChannelST `json:"channels,omitempty" groups:"api,config"`
}

type StorageST struct {
	Streams         map[string]StreamST `json:"streams,omitempty" groups:"api,config"`
	ChannelDefaults ChannelST           `json:"channel_defaults,omitempty" groups:"api,config"`
}

// Command line flag global variables
var debug bool

var (
	configFile string
	Storage    = NewStreamCore()
)

func NewStreamCore() *StorageST {

	flag.BoolVar(&debug, "debug", true, "set debug mode")
	flag.StringVar(&configFile, "config", "config/config.json", "config patch (/etc/server/config.json or config.json)")
	flag.Parse()
	var tmp StorageST
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		utils.Logger.Error(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		utils.Logger.Error(err.Error())
		os.Exit(1)
	}
	// debug = tmp.Server.Debug
	for i, i2 := range tmp.Streams {
		for i3, i4 := range i2.Channels {
			channel := tmp.ChannelDefaults
			err = mergo.Merge(&channel, i4)
			if err != nil {
				utils.Logger.Error(err.Error())
				os.Exit(1)
			}
			channel.clients = make(map[string]ClientST)
			channel.ack = time.Now().Add(-255 * time.Hour)
			// channel.hlsSegmentBuffer = make(map[int]SegmentOld)
			channel.signals = make(chan int, 100)
			i2.Channels[i3] = channel
		}
		tmp.Streams[i] = i2
	}
	return &tmp
}
