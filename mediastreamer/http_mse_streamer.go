package mediastreamer

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/vtpl1/go_rtsp_to_web/datamodel"
	"github.com/vtpl1/go_rtsp_to_web/utils"
	"github.com/vtpl1/go_rtsp_to_web/webstreamer"
	"github.com/vtpl1/vdk/format/mp4f"
)

// HTTPAPIServerStreamMSE func
func HTTPAPIServerStreamMSE(c *gin.Context) {
	conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
	if err != nil {
		return
	}

	defer func() {
		err = conn.Close()
		utils.Logger.Errorln(err)
	}()
	if !datamodel.Storage.StreamChannelExist(c.Param("uuid"), c.Param("channel")) {
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}

	if !webstreamer.RemoteAuthorization("WS", c.Param("uuid"), c.Param("channel"), c.Param("token"), c.ClientIP()) {
		utils.Logger.Errorln(webstreamer.ErrorStreamNotFound.Error())
		return
	}

	datamodel.Storage.StreamChannelRun(c.Param("uuid"), c.Param("channel"))
	err = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	cid, ch, _, err := datamodel.Storage.ClientAdd(c.Param("uuid"), c.Param("channel"), datamodel.MSE)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	defer datamodel.Storage.ClientDelete(c.Param("uuid"), cid, c.Param("channel"))
	codecs, err := datamodel.Storage.StreamChannelCodecs(c.Param("uuid"), c.Param("channel"))
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	muxerMSE := mp4f.NewMuxer(nil)
	err = muxerMSE.WriteHeader(codecs)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	meta, init := muxerMSE.GetInit(codecs)
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, append([]byte{9}, meta...))
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	err = wsutil.WriteServerMessage(conn, ws.OpBinary, init)
	if err != nil {
		utils.Logger.Errorln(err.Error())
		return
	}
	var videoStart bool
	controlExit := make(chan bool, 10)
	noClient := time.NewTimer(10 * time.Second)
	go func() {
		defer func() {
			controlExit <- true
		}()
		for {
			header, _, err := wsutil.NextReader(conn, ws.StateServerSide)
			if err != nil {
				utils.Logger.Errorln(err.Error())
				return
			}
			switch header.OpCode {
			case ws.OpPong:
				noClient.Reset(10 * time.Second)
			case ws.OpClose:
				return
			}
		}
	}()
	noVideo := time.NewTimer(10 * time.Second)
	pingTicker := time.NewTicker(500 * time.Millisecond)
	defer pingTicker.Stop()
	defer log.Println("client exit")
	for {
		select {

		case <-pingTicker.C:
			err = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			if err != nil {
				return
			}
			buf, err := ws.CompileFrame(ws.NewPingFrame(nil))
			if err != nil {
				return
			}
			_, err = conn.Write(buf)
			if err != nil {
				return
			}
		case <-controlExit:
			utils.Logger.Errorln("Client Reader Exit")
			return
		case <-noClient.C:
			utils.Logger.Errorln("Client OffLine Exit")
			return
		case <-noVideo.C:
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
			ready, buf, err := muxerMSE.WritePacket(*pck, false)
			if err != nil {
				utils.Logger.Errorln(err.Error())
				return
			}
			if ready {
				err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err != nil {
					utils.Logger.Errorln(err.Error())
					return
				}
				// err = websocket.Message.Send(ws, buf)
				err = wsutil.WriteServerMessage(conn, ws.OpBinary, buf)
				if err != nil {
					utils.Logger.Errorln(err.Error())
					return
				}
			}
		}
	}
}
