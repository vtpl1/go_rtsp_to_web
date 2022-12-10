package datamodel

import (
	"context"
	"encoding/json"
	"flag"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/imdario/mergo"
	"github.com/liip/sheriff"
	"github.com/vtpl1/go_rtsp_to_web/utils"
	"github.com/vtpl1/vdk/av"
)

func NewStreamCore() *StorageST {
	// flag.BoolVar(&debug, "debug", true, "set debug mode")
	flag.StringVar(&configFile, "config", "config/config.json", "config patch (/etc/server/config.json or config.json)")
	flag.Parse()
	var tmp StorageST
	tmp.mutex = &sync.RWMutex{}
	data, err := os.ReadFile(configFile)
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
	for i, stream := range tmp.Streams {
		stream.mutex = &sync.RWMutex{}
		for i3, i4 := range stream.Channels {
			channel := tmp.ChannelDefaults
			err = mergo.Merge(&channel, i4)
			if err != nil {
				utils.Logger.Error(err.Error())
				os.Exit(1)
			}
			channel.clients = make(map[string]ClientST)
			channel.ack = time.Now().Add(-255 * time.Hour)
			channel.hlsSegmentBuffer = make(map[int]SegmentOld)
			channel.hlsSegmentNumber = 0
			channel.signals = make(chan int, 100)
			stream.Channels[i3] = channel
		}
		tmp.Streams[i] = stream
	}
	return &tmp
}

// StreamChannelUnlock unlock status to no lock
func (obj *StorageST) StreamChannelUnlock(streamID string, channelID string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[streamID]; ok {
		if channelTmp, ok := streamTmp.Channels[channelID]; ok {
			channelTmp.runLock = false
			streamTmp.Channels[channelID] = channelTmp
			obj.Streams[streamID] = streamTmp
		}
	}
}

// StreamChannelStatus change stream status
func (obj *StorageST) StreamChannelStatus(key string, channelID string, val int) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.Status = val
			tmp.Channels[channelID] = channelTmp
			obj.Streams[key] = tmp
		}
	}
}

// StreamHLSFlush delete hls cache
func (obj *StorageST) StreamHLSFlush(uuid string, channelID string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.hlsSegmentBuffer = make(map[int]SegmentOld)
			channelTmp.hlsSegmentNumber = 0
			tmp.Channels[channelID] = channelTmp
			obj.Streams[uuid] = tmp
		}
	}
}

// StreamChannelCodecsUpdate update stream codec storage
func (obj *StorageST) StreamChannelCodecsUpdate(streamID string, channelID string, val []av.CodecData, sdp []byte) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[streamID]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.codecs = val
			channelTmp.sdp = sdp
			tmp.Channels[channelID] = channelTmp
			obj.Streams[streamID] = tmp
		}
	}
}

// NewHLSMuxer new muxer init
func (obj *StorageST) NewHLSMuxer(uuid string, channelID string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.hlsMuxer = NewHLSMuxer(uuid)
			tmp.Channels[channelID] = channelTmp
			obj.Streams[uuid] = tmp
		}
	}
}

// HLSMuxerClose close muxer
func (obj *StorageST) HLSMuxerClose(uuid string, channelID string) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.hlsMuxer.Close()
		}
	}
}

// ClientHas check is client ext
func (obj *StorageST) ClientHas(streamID string, channelID string) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return false
	}
	channelTmp, ok := streamTmp.Channels[channelID]
	if !ok {
		return false
	}
	if time.Since(channelTmp.ack).Seconds() > 30 {
		return false
	}
	return true
}

// StreamChannelCastProxy broadcast stream
func (obj *StorageST) StreamChannelCastProxy(key string, channelID string, val *[]byte) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			if len(channelTmp.clients) > 0 {
				for _, i2 := range channelTmp.clients {
					if i2.mode != RTSP {
						continue
					}
					if len(i2.outgoingRTPPacket) < 1000 {
						i2.outgoingRTPPacket <- val
					} else if len(i2.signals) < 10 {
						// send stop signals to client
						i2.signals <- SignalStreamStop
						err := i2.socket.Close()
						if err != nil {
							utils.Logger.Errorw(err.Error(), "stream", key, "channel", channelID)
						}
					}
				}
				channelTmp.ack = time.Now()
				tmp.Channels[channelID] = channelTmp
				obj.Streams[key] = tmp
			}
		}
	}
}

// StreamHLSAdd add hls seq to buffer
func (obj *StorageST) StreamHLSAdd(uuid string, channelID string, val []*av.Packet, dur time.Duration) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.hlsSegmentNumber++
			channelTmp.hlsSegmentBuffer[channelTmp.hlsSegmentNumber] = SegmentOld{data: val, dur: dur}
			if len(channelTmp.hlsSegmentBuffer) >= 6 {
				delete(channelTmp.hlsSegmentBuffer, channelTmp.hlsSegmentNumber-6-1)
			}
			tmp.Channels[channelID] = channelTmp
			obj.Streams[uuid] = tmp
		}
	}
}

// StreamChannelCast broadcast stream
func (obj *StorageST) StreamChannelCast(key string, channelID string, val *av.Packet) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[key]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			if len(channelTmp.clients) > 0 {
				for _, i2 := range channelTmp.clients {
					if i2.mode == RTSP {
						continue
					}
					if len(i2.outgoingAVPacket) < 1000 {
						i2.outgoingAVPacket <- val
					} else if len(i2.signals) < 10 {
						// send stop signals to client
						i2.signals <- SignalStreamStop
						//No need close socket only send signal to reader / writer socket closed if client go to offline
						/*
							err := i2.socket.Close()
							if err != nil {
								log.WithFields(logrus.Fields{
									"module":  "storage",
									"stream":  key,
									"channel": key,
									"func":    "CastProxy",
									"call":    "Close",
								}).Errorln(err.Error())
							}
						*/
					}
				}
				channelTmp.ack = time.Now()
				tmp.Channels[channelID] = channelTmp
				obj.Streams[key] = tmp
			}
		}
	}
}

// HlsMuxerSetFPS write packet
func (obj *StorageST) HlsMuxerSetFPS(uuid string, channelID string, fps int) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok && channelTmp.hlsMuxer != nil {
			channelTmp.hlsMuxer.SetFPS(fps)
			tmp.Channels[channelID] = channelTmp
			obj.Streams[uuid] = tmp
		}
	}
}

// HlsMuxerWritePacket write packet
func (obj *StorageST) HlsMuxerWritePacket(uuid string, channelID string, packet *av.Packet) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok && channelTmp.hlsMuxer != nil {
			channelTmp.hlsMuxer.WritePacket(packet)
		}
	}
}

// ServerHTTPDebug read debug options
func (obj *StorageST) ServerHTTPDebug() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDebug
}

// ServerRTSPPort read HTTP Port options
func (obj *StorageST) ServerRTSPPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.RTSPPort
}

// ServerHTTPPort read HTTP Port options
func (obj *StorageST) ServerHTTPPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPort
}

// ServerHTTPLogin read Login options
func (obj *StorageST) ServerHTTPLogin() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPLogin
}

// ServerHTTPPassword read Password options
func (obj *StorageST) ServerHTTPPassword() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPPassword
}

// ServerHTTPDemo read demo options
func (obj *StorageST) ServerHTTPDemo() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPDemo
}

// ServerICEServers read ICE servers
func (obj *StorageST) ServerICEServers() []string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.ICEServers
}

// ServerICEServers read ICE username
func (obj *StorageST) ServerICEUsername() string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.ICEUsername
}

// ServerICEServers read ICE credential
func (obj *StorageST) ServerICECredential() string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.ICECredential
}

// ServerWebRTCPortMin read WebRTC Port Min
func (obj *StorageST) ServerWebRTCPortMin() uint16 {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.WebRTCPortMin
}

// ServerWebRTCPortMax read WebRTC Port Max
func (obj *StorageST) ServerWebRTCPortMax() uint16 {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.Server.WebRTCPortMax
}

// Default www static file dir
var DefaultHTTPDir = "web"

// ServerHTTPDir
func (obj *StorageST) ServerHTTPDir() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if filepath.Clean(obj.Server.HTTPDir) == "." {
		return DefaultHTTPDir
	}
	return filepath.Clean(obj.Server.HTTPDir)
}

// ServerHTTPS read HTTPS Port options
func (obj *StorageST) ServerHTTPS() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPS
}

// ServerHTTPSAutoTLSEnable read HTTPS Port options
func (obj *StorageST) ServerHTTPSAutoTLSEnable() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSAutoTLSEnable
}

// ServerHTTPSAutoTLSName read HTTPS Port options
func (obj *StorageST) ServerHTTPSAutoTLSName() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSAutoTLSName
}

// ServerHTTPSCert read HTTPS Cert options
func (obj *StorageST) ServerHTTPSCert() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSCert
}

// ServerHTTPSKey read HTTPS Key options
func (obj *StorageST) ServerHTTPSKey() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSKey
}

// ServerHTTPSPort read HTTPS Port options
func (obj *StorageST) ServerHTTPSPort() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.HTTPSPort
}

// StreamsList list all stream
func (obj *StorageST) StreamsList() map[string]StreamST {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	tmp := make(map[string]StreamST)
	for i, i2 := range obj.Streams {
		tmp[i] = i2
	}
	return tmp
}

// StreamAdd add stream
func (obj *StorageST) StreamAdd(uuid string, val StreamST) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	// TODO create empty map bug save https://github.com/liip/sheriff empty not nil map[] != {} json
	// data, err := sheriff.Marshal(&sheriff.Options{
	//		Groups:     []string{"config"},
	//		ApiVersion: v2,
	//	}, obj)
	// Not Work map[] != {}
	if obj.Streams == nil {
		obj.Streams = make(map[string]StreamST)
	}
	if _, ok := obj.Streams[uuid]; ok {
		return ErrorStreamAlreadyExists
	}
	for i, i2 := range val.Channels {
		i2 = obj.StreamChannelMake(i2)
		if !i2.OnDemand {
			i2.runLock = true
			val.Channels[i] = i2
			go StreamServerRunStreamDo(uuid, i)
		} else {
			val.Channels[i] = i2
		}
	}
	obj.Streams[uuid] = val
	err := obj.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

func (obj *StorageST) SaveConfig() error {
	utils.Logger.Debugln("Saving configuration to", configFile)
	v2, err := version.NewVersion("2.0.0")
	if err != nil {
		return err
	}
	data, err := sheriff.Marshal(&sheriff.Options{
		Groups:     []string{"config"},
		ApiVersion: v2,
	}, obj)
	if err != nil {
		return err
	}
	res, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(configFile, res, 0o644)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return err
	}
	return nil
}

// StreamChannelMake check stream exist
func (obj *StorageST) StreamChannelMake(val ChannelST) ChannelST {
	channel := obj.ChannelDefaults
	if err := mergo.Merge(&channel, val); err != nil {
		// Just ignore the default values and continue
		channel = val
		utils.Logger.Errorln(err.Error())
	}
	// make client's
	channel.clients = make(map[string]ClientST)
	// make last ack
	channel.ack = time.Now().Add(-255 * time.Hour)
	// make hls buffer
	channel.hlsSegmentBuffer = make(map[int]SegmentOld)
	// make signals buffer chain
	channel.signals = make(chan int, 100)
	return channel
}

// StreamEdit edit stream
func (obj *StorageST) StreamEdit(uuid string, val StreamST) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		for i, i2 := range tmp.Channels {
			if i2.runLock {
				tmp.Channels[i] = i2
				obj.Streams[uuid] = tmp
				i2.signals <- SignalStreamStop
			}
		}
		for i3, i4 := range val.Channels {
			i4 = obj.StreamChannelMake(i4)
			if !i4.OnDemand {
				i4.runLock = true
				val.Channels[i3] = i4
				go StreamServerRunStreamDo(uuid, i3)
			} else {
				val.Channels[i3] = i4
			}
		}
		obj.Streams[uuid] = val
		err := obj.SaveConfig()
		if err != nil {
			return err
		}
		return nil
	}
	return ErrorStreamNotFound
}

// StreamDelete stream
func (obj *StorageST) StreamDelete(uuid string) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		for _, i2 := range tmp.Channels {
			if i2.runLock {
				i2.signals <- SignalStreamStop
			}
		}
		delete(obj.Streams, uuid)
		err := obj.SaveConfig()
		if err != nil {
			return err
		}
		return nil
	}
	return ErrorStreamNotFound
}

// StreamReload reload stream
func (obj *StorageST) StreamReload(uuid string) error {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		for _, i2 := range tmp.Channels {
			if i2.runLock {
				i2.signals <- SignalStreamRestart
			}
		}
		return nil
	}
	return ErrorStreamNotFound
}

// StreamInfo return stream info
func (obj *StorageST) StreamInfo(uuid string) (*StreamST, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		return &tmp, nil
	}
	return nil, ErrorStreamNotFound
}

// StreamChannelAdd add stream
func (obj *StorageST) StreamChannelAdd(uuid string, channelID string, val ChannelST) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[uuid]; !ok {
		return ErrorStreamNotFound
	}
	if _, ok := obj.Streams[uuid].Channels[channelID]; ok {
		return ErrorStreamChannelAlreadyExists
	}
	val = obj.StreamChannelMake(val)
	obj.Streams[uuid].Channels[channelID] = val
	if !val.OnDemand {
		val.runLock = true
		go StreamServerRunStreamDo(uuid, channelID)
	}
	err := obj.SaveConfig()
	if err != nil {
		return err
	}
	return nil
}

// StreamChannelDelete stream
func (obj *StorageST) StreamChannelDelete(uuid string, channelID string) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			if channelTmp.runLock {
				channelTmp.signals <- SignalStreamStop
			}
			delete(obj.Streams[uuid].Channels, channelID)
			err := obj.SaveConfig()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrorStreamNotFound
}

// StreamEdit edit stream
func (obj *StorageST) StreamChannelEdit(uuid string, channelID string, val ChannelST) error {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if currentChannel, ok := tmp.Channels[channelID]; ok {
			if currentChannel.runLock {
				currentChannel.signals <- SignalStreamStop
			}
			val = obj.StreamChannelMake(val)
			obj.Streams[uuid].Channels[channelID] = val
			if !val.OnDemand {
				val.runLock = true
				go StreamServerRunStreamDo(uuid, channelID)
			}
			err := obj.SaveConfig()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return ErrorStreamNotFound
}

// StreamChannelExist check stream exist
func (obj *StorageST) StreamChannelExist(streamID string, channelID string) bool {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[streamID]; ok {
		if channelTmp, ok := streamTmp.Channels[channelID]; ok {
			channelTmp.ack = time.Now()
			streamTmp.Channels[channelID] = channelTmp
			obj.Streams[streamID] = streamTmp
			return ok
		}
	}
	return false
}

// StreamChannelCodecs get stream codec storage or wait
func (obj *StorageST) StreamChannelCodecs(streamID string, channelID string) ([]av.CodecData, error) {
	for i := 0; i < 100; i++ {
		obj.mutex.RLock()
		tmp, ok := obj.Streams[streamID]
		obj.mutex.RUnlock()
		if !ok {
			return nil, ErrorStreamNotFound
		}
		channelTmp, ok := tmp.Channels[channelID]
		if !ok {
			return nil, ErrorStreamChannelNotFound
		}

		if channelTmp.codecs != nil {
			return channelTmp.codecs, nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil, ErrorStreamChannelCodecNotFound
}

// StreamInfo return stream info
func (obj *StorageST) StreamChannelInfo(uuid string, channelID string) (*ChannelST, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			return &channelTmp, nil
		}
	}
	return nil, ErrorStreamNotFound
}

// StreamChannelReload reload stream
func (obj *StorageST) StreamChannelReload(uuid string, channelID string) error {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			channelTmp.signals <- SignalStreamRestart
			return nil
		}
	}
	return ErrorStreamNotFound
}

// NewHLSMuxer Segments
func NewHLSMuxer(uuid string) *MuxerHLS {
	ctx, cancel := context.WithCancel(context.Background())
	return &MuxerHLS{
		UUID:           uuid,
		MSN:            -1,
		Segments:       make(map[int]*Segment),
		FragmentCtx:    ctx,
		FragmentCancel: cancel,
	}
}

// SetFPS func
func (element *MuxerHLS) SetFPS(fps int) {
	element.FPS = fps
}

// WritePacket func
func (element *MuxerHLS) WritePacket(packet *av.Packet) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	// TODO delete packet.IsKeyFrame if need no EXT-X-INDEPENDENT-SEGMENTS

	if !packet.IsKeyFrame && element.CurrentSegment == nil {
		// Wait for the first keyframe before initializing
		return
	}
	if packet.IsKeyFrame && (element.CurrentSegment == nil || element.CurrentSegment.GetDuration().Seconds() >= 4) {
		if element.CurrentSegment != nil {
			element.CurrentSegment.Close()
			if len(element.Segments) > 6 {
				delete(element.Segments, element.MSN-6)
				element.MediaSequence++
			}
		}
		element.CurrentSegment = element.NewSegment()
		element.CurrentSegment.SetFPS(element.FPS)
	}
	element.CurrentSegment.WritePacket(packet)
	CurrentFragmentID := element.CurrentSegment.GetFragmentID()
	if CurrentFragmentID != element.CurrentFragmentID {
		element.UpdateIndexM3u8()
	}
	element.CurrentFragmentID = CurrentFragmentID
}

// NewSegment func
func (element *MuxerHLS) NewSegment() *Segment {
	res := &Segment{
		Fragment:          make(map[int]*Fragment),
		CurrentFragmentID: -1, // Default fragment -1
	}
	// Increase MSN
	element.MSN++
	element.Segments[element.MSN] = res
	return res
}

// UpdateIndexM3u8 func
func (element *MuxerHLS) UpdateIndexM3u8() {
	var header string
	var body string
	var partTarget time.Duration
	var segmentTarget time.Duration
	segmentTarget = time.Second * 2
	for _, segmentKey := range element.SortSegments(element.Segments) {
		for _, fragmentKey := range element.SortFragment(element.Segments[segmentKey].Fragment) {
			if element.Segments[segmentKey].Fragment[fragmentKey].Finish {
				var independent string
				if element.Segments[segmentKey].Fragment[fragmentKey].Independent {
					independent = ",INDEPENDENT=YES"
				}
				body += "#EXT-X-PART:DURATION=" + strconv.FormatFloat(element.Segments[segmentKey].Fragment[fragmentKey].GetDuration().Seconds(), 'f', 5, 64) + "" + independent + ",URI=\"fragment/" + strconv.Itoa(segmentKey) + "/" + strconv.Itoa(fragmentKey) + "/0qrm9ru6." + strconv.Itoa(fragmentKey) + ".m4s\"\n"
				partTarget = element.Segments[segmentKey].Fragment[fragmentKey].Duration
			} else {
				body += "#EXT-X-PRELOAD-HINT:TYPE=PART,URI=\"fragment/" + strconv.Itoa(segmentKey) + "/" + strconv.Itoa(fragmentKey) + "/0qrm9ru6." + strconv.Itoa(fragmentKey) + ".m4s\"\n"
			}
		}
		if element.Segments[segmentKey].Finish {
			segmentTarget = element.Segments[segmentKey].Duration
			body += "#EXT-X-PROGRAM-DATE-TIME:" + element.Segments[segmentKey].Time.Format("2006-01-02T15:04:05.000000Z") + "\n#EXTINF:" + strconv.FormatFloat(element.Segments[segmentKey].Duration.Seconds(), 'f', 5, 64) + ",\n"
			body += "segment/" + strconv.Itoa(segmentKey) + "/" + element.UUID + "." + strconv.Itoa(segmentKey) + ".m4s\n"
		}
	}
	header += "#EXTM3U\n"
	header += "#EXT-X-TARGETDURATION:" + strconv.Itoa(int(math.Round(segmentTarget.Seconds()))) + "\n"
	header += "#EXT-X-VERSION:7\n"
	header += "#EXT-X-INDEPENDENT-SEGMENTS\n"
	header += "#EXT-X-SERVER-CONTROL:CAN-BLOCK-RELOAD=YES,PART-HOLD-BACK=" + strconv.FormatFloat(partTarget.Seconds()*4, 'f', 5, 64) + ",HOLD-BACK=" + strconv.FormatFloat(segmentTarget.Seconds()*4, 'f', 5, 64) + "\n"
	header += "#EXT-X-MAP:URI=\"init.mp4\"\n"
	header += "#EXT-X-PART-INF:PART-TARGET=" + strconv.FormatFloat(partTarget.Seconds(), 'f', 5, 64) + "\n"
	header += "#EXT-X-MEDIA-SEQUENCE:" + strconv.Itoa(element.MediaSequence) + "\n"
	header += body
	element.CacheM3U8 = header
	element.PlaylistUpdate()
}

// SortSegments fuc
func (element *MuxerHLS) SortSegments(val map[int]*Segment) []int {
	keys := make([]int, len(val))
	i := 0
	for k := range val {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

// SortFragment func
func (element *MuxerHLS) SortFragment(val map[int]*Fragment) []int {
	keys := make([]int, len(val))
	i := 0
	for k := range val {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	return keys
}

// PlaylistUpdate func
func (element *MuxerHLS) PlaylistUpdate() {
	element.FragmentCancel()
	element.FragmentCtx, element.FragmentCancel = context.WithCancel(context.Background())
}

func (element *MuxerHLS) Close() {
}

// GetIndexM3u8 func
func (element *MuxerHLS) GetIndexM3u8(needMSN int, needPart int) (string, error) {
	element.mutex.Lock()
	if len(element.CacheM3U8) != 0 && ((needMSN == -1 || needPart == -1) || (needMSN-element.MSN > 1) || (needMSN == element.MSN && needPart < element.CurrentFragmentID)) {
		element.mutex.Unlock()
		return element.CacheM3U8, nil
	} else {
		element.mutex.Unlock()
		index, err := element.WaitIndex(time.Second*3, needMSN, needPart)
		if err != nil {
			return "", err
		}
		return index, err
	}
}

// WaitIndex func
func (element *MuxerHLS) WaitIndex(timeOut time.Duration, segment, fragment int) (string, error) {
	for {
		select {
		case <-time.After(timeOut):
			return "", ErrorStreamNotFound
		case <-element.FragmentCtx.Done():
			element.mutex.Lock()
			if element.MSN < segment || (element.MSN == segment && element.CurrentFragmentID < fragment) {
				utils.Logger.Infoln("wait req", element.MSN, element.CurrentFragmentID, segment, fragment)
				element.mutex.Unlock()
				continue
			}
			element.mutex.Unlock()
			return element.CacheM3U8, nil
		}
	}
}

// GetSegment func
func (element *MuxerHLS) GetSegment(segment int) ([]*av.Packet, error) {
	element.mutex.Lock()
	defer element.mutex.Unlock()
	if segmentTmp, ok := element.Segments[segment]; ok && len(segmentTmp.Fragment) > 0 {
		var res []*av.Packet
		for _, v := range element.SortFragment(segmentTmp.Fragment) {
			res = append(res, segmentTmp.Fragment[v].Packets...)
		}
		return res, nil
	}
	return nil, ErrorStreamNotFound
}

// GetFragment func
func (element *MuxerHLS) GetFragment(segment int, fragment int) ([]*av.Packet, error) {
	element.mutex.Lock()
	if segmentTmp, segmentTmpOK := element.Segments[segment]; segmentTmpOK {
		if fragmentTmp, fragmentTmpOK := segmentTmp.Fragment[fragment]; fragmentTmpOK {
			if fragmentTmp.Finish {
				element.mutex.Unlock()
				return fragmentTmp.Packets, nil
			} else {
				element.mutex.Unlock()
				pck, err := element.WaitFragment(time.Second*1, segment, fragment)
				if err != nil {
					return nil, err
				}
				return pck, err
			}
		}
	}
	element.mutex.Unlock()
	return nil, ErrorStreamNotFound
}

// WaitFragment func
func (element *MuxerHLS) WaitFragment(timeOut time.Duration, segment, fragment int) ([]*av.Packet, error) {
	select {
	case <-time.After(timeOut):
		return nil, ErrorStreamNotFound
	case <-element.FragmentCtx.Done():
		element.mutex.Lock()
		defer element.mutex.Unlock()
		if segmentTmp, segmentTmpOK := element.Segments[segment]; segmentTmpOK {
			if fragmentTmp, fragmentTmpOK := segmentTmp.Fragment[fragment]; fragmentTmpOK {
				if fragmentTmp.Finish {
					return fragmentTmp.Packets, nil
				}
			}
		}
		return nil, ErrorStreamNotFound
	}
}

// GetDuration func
func (element *Segment) GetDuration() time.Duration {
	return element.Duration
}

// SetFPS func
func (element *Segment) SetFPS(fps int) {
	element.FPS = fps
}

// WritePacket func
func (element *Segment) WritePacket(packet *av.Packet) {
	if element.CurrentFragment == nil || element.CurrentFragment.GetDuration().Milliseconds() >= element.FragmentMS(element.FPS) {
		if element.CurrentFragment != nil {
			element.CurrentFragment.Close()
		}
		element.CurrentFragmentID++
		element.CurrentFragment = element.NewFragment()
	}
	element.Duration += packet.Duration
	element.CurrentFragment.WritePacket(packet)
}

// FragmentMS func
func (element *Segment) FragmentMS(fps int) int64 {
	for i := 6; i >= 1; i-- {
		if fps%i == 0 {
			return int64(float64(1000) / float64(fps) * float64(i))
		}
	}
	return 100
}

// Close segment func
func (element *Segment) Close() {
	element.Finish = true
	if element.CurrentFragment != nil {
		element.CurrentFragment.Close()
	}
}

// GetFragmentID func
func (element *Segment) GetFragmentID() int {
	return element.CurrentFragmentID
}

// NewFragment open new fragment
func (element *Segment) NewFragment() *Fragment {
	res := &Fragment{}
	element.Fragment[element.CurrentFragmentID] = res
	return res
}

// GetDuration return fragment dur
func (element *Fragment) GetDuration() time.Duration {
	return element.Duration
}

// WritePacket to fragment func
func (element *Fragment) WritePacket(packet *av.Packet) {
	// increase fragment dur
	element.Duration += packet.Duration
	// Independent if have key
	if packet.IsKeyFrame {
		element.Independent = true
	}
	// append packet to slice of packet
	element.Packets = append(element.Packets, packet)
}

// Close fragment block func
func (element *Fragment) Close() {
	// TODO add callback func
	// finalize fragment
	element.Finish = true
}
