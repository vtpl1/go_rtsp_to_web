package datamodel

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/imdario/mergo"
	"github.com/vtpl1/go_rtsp_to_web/utils"
	"github.com/vtpl1/vdk/av"
)

func NewStreamCore() *StorageST {
	// flag.BoolVar(&debug, "debug", true, "set debug mode")
	flag.StringVar(&configFile, "config", "config/config.json", "config patch (/etc/server/config.json or config.json)")
	flag.Parse()
	var tmp StorageST
	tmp.mutex = &sync.RWMutex{}
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
	if time.Now().Sub(channelTmp.ack).Seconds() > 30 {
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

var (
	//Default www static file dir
	DefaultHTTPDir = "web"
)

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
