package mediastreamer

import (
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vtpl1/go_rtsp_to_web/datamodel"
	"github.com/vtpl1/go_rtsp_to_web/utils"
	"github.com/vtpl1/go_rtsp_to_web/webstreamer"
	"github.com/vtpl1/vdk"
	"github.com/vtpl1/vdk/format/ts"
)

// HTTPAPIServerStreamHLSM3U8 send client m3u8 play list
func HTTPAPIServerStreamHLSM3U8(c *gin.Context) {
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
	// If stream mode on_demand need wait ready segment's
	for i := 0; i < 40; i++ {
		index, seq, err := datamodel.Storage.StreamHLSm3u8(c.Param("uuid"), c.Param("channel"))
		if err != nil {
			c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
			utils.Logger.Errorln(err.Error())
			return
		}
		if seq >= 6 {
			_, err := c.Writer.Write([]byte(index))
			if err != nil {
				c.IndentedJSON(400, webstreamer.Message{Status: 0, Payload: err.Error()})
				utils.Logger.Errorln(err.Error())
				return
			}
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// HTTPAPIServerStreamHLSTS send client ts segment
func HTTPAPIServerStreamHLSTS(c *gin.Context) {

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
	outfile := bytes.NewBuffer([]byte{})
	Muxer := ts.NewMuxer(outfile)
	Muxer.PaddingToMakeCounterCont = true
	err = Muxer.WriteHeader(codecs)
	if err != nil {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	seqData, err := datamodel.Storage.StreamHLSTS(c.Param("uuid"), c.Param("channel"), vdk.StringToInt(c.Param("seq")))
	if err != nil {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	if len(seqData) == 0 {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: webstreamer.ErrorStreamNotHLSSegments.Error()})
		utils.Logger.Errorln(webstreamer.ErrorStreamNotHLSSegments.Error())
		return
	}
	for _, v := range seqData {
		v.CompositionTime = 1
		err = Muxer.WritePacket(*v)
		if err != nil {
			c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
			utils.Logger.Errorln(err.Error())
			return
		}
	}
	err = Muxer.WriteTrailer()
	if err != nil {
		c.IndentedJSON(500, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
	_, err = c.Writer.Write(outfile.Bytes())
	if err != nil {
		c.IndentedJSON(400, webstreamer.Message{Status: 0, Payload: err.Error()})
		utils.Logger.Errorln(err.Error())
		return
	}
}
