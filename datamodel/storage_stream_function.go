package datamodel

import (
	"sort"
	"strconv"
	"time"

	"github.com/vtpl1/vdk/av"
)

// ServerTokenEnable read HTTPS Key options
func (obj *StorageST) ServerTokenEnable() bool {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.Token.Enable
}

// ServerTokenBackend read HTTPS Key options
func (obj *StorageST) ServerTokenBackend() string {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	return obj.Server.Token.Backend
}

// StreamChannelRun one stream and lock
func (obj *StorageST) StreamChannelRun(streamID string, channelID string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if streamTmp, ok := obj.Streams[streamID]; ok {
		if channelTmp, ok := streamTmp.Channels[channelID]; ok {
			if !channelTmp.runLock {
				channelTmp.runLock = true
				streamTmp.Channels[channelID] = channelTmp
				obj.Streams[streamID] = streamTmp
				go StreamServerRunStreamDo(streamID, channelID)
			}
		}
	}
}

// StreamHLSm3u8 get hls m3u8 list
func (obj *StorageST) StreamHLSm3u8(uuid string, channelID string) (string, int, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			var out string
			// TODO fix  it
			out += "#EXTM3U\r\n#EXT-X-TARGETDURATION:4\r\n#EXT-X-VERSION:4\r\n#EXT-X-MEDIA-SEQUENCE:" + strconv.Itoa(channelTmp.hlsSegmentNumber) + "\r\n"
			var keys []int
			for k := range channelTmp.hlsSegmentBuffer {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			var count int
			for _, i := range keys {
				count++
				out += "#EXTINF:" + strconv.FormatFloat(channelTmp.hlsSegmentBuffer[i].dur.Seconds(), 'f', 1, 64) + ",\r\nsegment/" + strconv.Itoa(i) + "/file.ts\r\n"

			}
			return out, count, nil
		}
	}
	return "", 0, ErrorStreamNotFound
}

// StreamHLSTS send hls segment buffer to clients
func (obj *StorageST) StreamHLSTS(uuid string, channelID string, seq int) ([]*av.Packet, error) {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			if tmp, ok := channelTmp.hlsSegmentBuffer[seq]; ok {
				return tmp.data, nil
			}
		}
	}
	return nil, ErrorStreamNotFound
}

// HLSMuxerM3U8 get m3u8 list
func (obj *StorageST) HLSMuxerM3U8(uuid string, channelID string, msn, part int) (string, error) {
	obj.mutex.Lock()
	tmp, ok := obj.Streams[uuid]
	obj.mutex.Unlock()
	if ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			index, err := channelTmp.hlsMuxer.GetIndexM3u8(msn, part)
			return index, err
		}
	}
	return "", ErrorStreamNotFound
}

// HLSMuxerSegment get segment
func (obj *StorageST) HLSMuxerSegment(uuid string, channelID string, segment int) ([]*av.Packet, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if tmp, ok := obj.Streams[uuid]; ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			return channelTmp.hlsMuxer.GetSegment(segment)
		}
	}
	return nil, ErrorStreamChannelNotFound
}

// HLSMuxerFragment get fragment
func (obj *StorageST) HLSMuxerFragment(uuid string, channelID string, segment, fragment int) ([]*av.Packet, error) {
	obj.mutex.Lock()
	tmp, ok := obj.Streams[uuid]
	obj.mutex.Unlock()
	if ok {
		if channelTmp, ok := tmp.Channels[channelID]; ok {
			packet, err := channelTmp.hlsMuxer.GetFragment(segment, fragment)
			return packet, err
		}
	}
	return nil, ErrorStreamChannelNotFound
}

// ClientAdd Add New Client to Translations
func (obj *StorageST) ClientAdd(streamID string, channelID string, mode int) (string, chan *av.Packet, chan *[]byte, error) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	streamTmp, ok := obj.Streams[streamID]
	if !ok {
		return "", nil, nil, ErrorStreamNotFound
	}
	// Generate UUID client
	cid, err := generateUUID()
	if err != nil {
		return "", nil, nil, err
	}
	chAV := make(chan *av.Packet, 2000)
	chRTP := make(chan *[]byte, 2000)
	channelTmp, ok := streamTmp.Channels[channelID]
	if !ok {
		return "", nil, nil, ErrorStreamNotFound
	}

	channelTmp.clients[cid] = ClientST{mode: mode, outgoingAVPacket: chAV, outgoingRTPPacket: chRTP, signals: make(chan int, 100)}
	channelTmp.ack = time.Now()
	streamTmp.Channels[channelID] = channelTmp
	obj.Streams[streamID] = streamTmp
	return cid, chAV, chRTP, nil
}

// ClientDelete Delete Client
func (obj *StorageST) ClientDelete(streamID string, cid string, channelID string) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if _, ok := obj.Streams[streamID]; ok {
		delete(obj.Streams[streamID].Channels[channelID].clients, cid)
	}
}

// StreamChannelSDP codec storage or wait
func (obj *StorageST) StreamChannelSDP(streamID string, channelID string) ([]byte, error) {
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

		if len(channelTmp.sdp) > 0 {
			return channelTmp.sdp, nil
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil, ErrorStreamNotFound
}
