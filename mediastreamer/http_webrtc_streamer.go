package mediastreamer

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vtpl1/go_rtsp_to_web/datamodel"
	"github.com/vtpl1/go_rtsp_to_web/utils"
	"github.com/vtpl1/go_rtsp_to_web/webstreamer"
	webrtc "github.com/vtpl1/vdk/format/webrtcv3"
)

// HTTPAPIServerStreamWebRTC stream video over WebRTC
func HTTPAPIServerStreamWebRTC(c *gin.Context) {
	if !datamodel.Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: webstreamer.ErrorStreamNotFound.Error()})
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}

	if !webstreamer.RemoteAuthorization("WebRTC", c.Param("uuid"), c.Param("channel"), c.Query("token"), c.ClientIP()) {
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}

	datamodel.Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
	codecs, err := datamodel.Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	muxerWebRTC := webrtc.NewMuxer(webrtc.Options{ICEServers: datamodel.Storage.ServerICEServers(), ICEUsername: datamodel.Storage.ServerICEUsername(), ICECredential: datamodel.Storage.ServerICECredential(), PortMin: datamodel.Storage.ServerWebRTCPortMin(), PortMax: datamodel.Storage.ServerWebRTCPortMax()})
	answer, err := muxerWebRTC.WriteHeader(codecs, c.PostForm("data"))
	if err != nil {
		c.IndentedJSON(400, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	_, err = c.Writer.Write([]byte(answer))
	if err != nil {
		c.IndentedJSON(400, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	go func() {
		cid, ch, _, err := datamodel.Storage.ClientAdd(c.Param("uuid"), c.Param("channel"), datamodel.WEBRTC)
		if err != nil {
			c.IndentedJSON(400, webstreamer.Message{Status: 0, Payload: err.Error()})
			utils.Logger.Errorln(err.Error())
			return
		}
		defer datamodel.Storage.ClientDelete(c.Param("uuid"), cid, c.Param("channel"))
		var videoStart bool
		noVideo := time.NewTimer(10 * time.Second)
		for {
			select {
			case <-noVideo.C:
				//				c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNoVideo.Error()})
				utils.Logger.Errorln(webstreamer.ErrorStreamNoVideo.Error())
				return
			case pck := <-ch:
				if pck.IsKeyFrame {
					noVideo.Reset(10 * time.Second)
					videoStart = true
				}
				if !videoStart {
					continue
				}
				err = muxerWebRTC.WritePacket(*pck)
				if err != nil {
					utils.Logger.Errorln(err.Error())
					return
				}
			}
		}
	}()
}
