package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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
	clients            map[string]ClientST
	ack                time.Time
	// hlsMuxer           *MuxerHLS `json:"-"`
}

// ServerST stream storage section
type StreamST struct {
	Name     string               `json:"name,omitempty" groups:"api,config"`
	Channels map[string]ChannelST `json:"channels,omitempty" groups:"api,config"`
}

type StorageST struct {
	Streams map[string]StreamST `json:"streams,omitempty" groups:"api,config"`
}

// Command line flag global variables
var debug bool
var configFile string
var Storage = NewStreamCore()

func NewStreamCore() *StorageST {
	logger := log.WithFields(logrus.Fields{
		"module": "config",
		"func":   "NewStreamCore",
		"call":   "ReadFile",
	})
	flag.BoolVar(&debug, "debug", true, "set debug mode")
	flag.StringVar(&configFile, "config", "config/config.json", "config patch (/etc/server/config.json or config.json)")
	flag.Parse()
	var tmp StorageST
	data, err := ioutil.ReadFile(configFile)
	if err == nil {
		logger.Errorln(err.Error())
		os.Exit(1)
	}
	err = json.Unmarshal(data, &tmp)
	return &tmp
}