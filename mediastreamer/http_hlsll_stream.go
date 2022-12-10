package mediastreamer

import (
	"github.com/gin-gonic/gin"
	"github.com/vtpl1/go_rtsp_to_web/datamodel"
	"github.com/vtpl1/go_rtsp_to_web/utils"
	"github.com/vtpl1/go_rtsp_to_web/webstreamer"
	"github.com/vtpl1/vdk"
	"github.com/vtpl1/vdk/format/mp4f"
)

// HTTPAPIServerStreamHLSLLM3U8 send client m3u8 play list
func HTTPAPIServerStreamHLSLLM3U8(c *gin.Context) {
	if !datamodel.Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: webstreamer.ErrorStreamNotFound.Error()})
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}
	c.Header("Content-Type", "application/x-mpegURL")
	datamodel.Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
	index, err := datamodel.Storage.HLSMuxerM3U8(c.Param("uuid"), c.Param("channel"), vdk.StringToInt(c.DefaultQuery("_HLS_msn", "-1")), vdk.StringToInt(c.DefaultQuery("_HLS_part", "-1")))
	if err != nil {
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}
	_, err = c.Writer.Write([]byte(index))
	if err != nil {
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}
}

// HTTPAPIServerStreamHLSLLInit send client ts segment
func HTTPAPIServerStreamHLSLLInit(c *gin.Context) {
	if !datamodel.Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: webstreamer.ErrorStreamNotFound.Error()})
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}

	if !webstreamer.RemoteAuthorization("HLS", c.Param("uuid"), c.Param("channel"), c.Param("token"), c.ClientIP()) {
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}

	c.Header("Content-Type", "application/x-mpegURL")
	datamodel.Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
	codecs, err := datamodel.Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	Muxer := mp4f.NewMuxer(nil)
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	c.Header("Content-Type", "video/mp4")
	_, buf := Muxer.GetInit(codecs)
	_, err = c.Writer.Write(buf)
	if err != nil {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
}

// HTTPAPIServerStreamHLSLLM4Segment send client ts segment
func HTTPAPIServerStreamHLSLLM4Segment(c *gin.Context) {
	c.Header("Content-Type", "video/mp4")
	if !datamodel.Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: webstreamer.ErrorStreamNotFound.Error()})
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}
	codecs, err := datamodel.Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	if codecs == nil {
		utils.Logger.Errorln("Codec Null")
		return
	}
	Muxer := mp4f.NewMuxer(nil)
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	seqData, err := datamodel.Storage.HLSMuxerSegment(c.Param("uuid"), c.Param("channel"), vdk.StringToInt(c.Param("segment")))
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	for _, v := range seqData {
		err = Muxer.WritePacket4(*v)
		if err != nil {
			utils.Logger.Errorln(err.Error())
			return
		}
	}
	buf := Muxer.Finalize()
	_, err = c.Writer.Write(buf)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
}

// HTTPAPIServerStreamHLSLLM4Fragment send client ts segment
func HTTPAPIServerStreamHLSLLM4Fragment(c *gin.Context) {
	c.Header("Content-Type", "video/mp4")
	if !datamodel.Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: webstreamer.ErrorStreamNotFound.Error()})
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}
	codecs, err := datamodel.Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	if codecs == nil {
		utils.Logger.Errorln("Codec Null")
		return
	}
	Muxer := mp4f.NewMuxer(nil)
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	seqData, err := datamodel.Storage.HLSMuxerFragment(c.Param("uuid"), c.Param("channel"), vdk.StringToInt(c.Param("segment")), vdk.StringToInt(c.Param("fragment")))
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	for _, v := range seqData {
		err = Muxer.WritePacket4(*v)
		if err != nil {
			utils.Logger.Errorln(err.Error())
			return
		}
	}
	buf := Muxer.Finalize()
	_, err = c.Writer.Write(buf)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
}
