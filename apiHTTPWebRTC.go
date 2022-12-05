package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	webrtc "github.com/vtpl1/vdk/format/webrtcv3"
)

// HTTPAPIServerStreamWebRTC stream video over WebRTC
func HTTPAPIServerStreamWebRTC(c *gin.Context) {
	requestLogger := log.WithFields(logrus.Fields{
		"module":  "http_webrtc",
		"stream":  c.Param("uuid"),
		"channel": c.Param("channel"),
		"func":    "HTTPAPIServerStreamWebRTC",
	})

	if !Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNotFound.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamChannelExist",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}

	if !RemoteAuthorization("WebRTC", c.Param("uuid"), c.Param("channel"), c.Query("token"), c.ClientIP()) {
		requestLogger.WithFields(logrus.Fields{
			"call": "RemoteAuthorization",
		}).Errorln(ErrorStreamNotFound.Error())
		return
	}

	Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
	codecs, err := Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "StreamCodecs",
		}).Errorln(err.Error())
		return
	}
	muxerWebRTC := webrtc.NewMuxer(webrtc.Options{ICEServers: Storage.ServerICEServers(), ICEUsername: Storage.ServerICEUsername(), ICECredential: Storage.ServerICECredential(), PortMin: Storage.ServerWebRTCPortMin(), PortMax: Storage.ServerWebRTCPortMax()})
	answer, err := muxerWebRTC.WriteHeader(codecs, c.PostForm("data"))
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "WriteHeader",
		}).Errorln(err.Error())
		return
	}
	_, err = c.Writer.Write([]byte(answer))
	if err != nil {
		c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
		requestLogger.WithFields(logrus.Fields{
			"call": "Write",
		}).Errorln(err.Error())
		return
	}
	go func() {
		cid, ch, _, err := Storage.ClientAdd(c.Param("uuid"), c.Param("channel"), WEBRTC)
		if err != nil {
			c.IndentedJSON(400, Message{Status: 0, Payload: err.Error()})
			requestLogger.WithFields(logrus.Fields{
				"call": "ClientAdd",
			}).Errorln(err.Error())
			return
		}
		defer Storage.ClientDelete(c.Param("uuid"), cid, c.Param("channel"))
		var videoStart bool
		noVideo := time.NewTimer(10 * time.Second)
		for {
			select {
			case <-noVideo.C:
				//				c.IndentedJSON(500, Message{Status: 0, Payload: ErrorStreamNoVideo.Error()})
				requestLogger.WithFields(logrus.Fields{
					"call": "ErrorStreamNoVideo",
				}).Errorln(ErrorStreamNoVideo.Error())
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
					requestLogger.WithFields(logrus.Fields{
						"call": "WritePacket",
					}).Errorln(err.Error())
					return
				}
			}
		}
	}()
}
