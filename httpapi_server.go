package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/vtpl1/go_rtsp_to_web/datamodel"
	"github.com/vtpl1/go_rtsp_to_web/mediastreamer"
	"github.com/vtpl1/go_rtsp_to_web/utils"
	"github.com/vtpl1/go_rtsp_to_web/webrenderer"
	"github.com/vtpl1/go_rtsp_to_web/webstreamer"
)

func HTTPAPIServer() {
	utils.Logger.Infof("HTTP Server start %s", datamodel.Storage.ServerHTTPPort())
	var public *gin.Engine
	if !datamodel.Storage.ServerHTTPDebug() {
		gin.SetMode(gin.ReleaseMode)
		public = gin.New()
	} else {
		gin.SetMode(gin.DebugMode)
		public = gin.Default()
	}
	public.Use(crossOrigin())
	// Add private login password protect methods
	private := public.Group("/")
	if datamodel.Storage.ServerHTTPLogin() != "" && datamodel.Storage.ServerHTTPPassword() != "" {
		private.Use(gin.BasicAuth(gin.Accounts{datamodel.Storage.ServerHTTPLogin(): datamodel.Storage.ServerHTTPPassword()}))
	}

	/*
		Static HTML Files Demo Mode
	*/

	if datamodel.Storage.ServerHTTPDemo() {
		public.LoadHTMLGlob(datamodel.Storage.ServerHTTPDir() + "/templates/*")
		public.GET("/", webrenderer.HTTPAPIServerIndex)
		public.GET("/pages/stream/list", webrenderer.HTTPAPIStreamList)
		public.GET("/pages/stream/add", webrenderer.HTTPAPIAddStream)
		public.GET("/pages/stream/edit/:uuid", webrenderer.HTTPAPIEditStream)
		public.GET("/pages/player/hls/:uuid/:channel", webrenderer.HTTPAPIPlayHls)
		public.GET("/pages/player/mse/:uuid/:channel", webrenderer.HTTPAPIPlayMse)
		public.GET("/pages/player/webrtc/:uuid/:channel", webrenderer.HTTPAPIPlayWebrtc)
		public.GET("/pages/multiview", webrenderer.HTTPAPIMultiview)
		public.Any("/pages/multiview/full", webrenderer.HTTPAPIFullScreenMultiView)
		public.GET("/pages/documentation", webrenderer.HTTPAPIServerDocumentation)
		public.GET("/pages/player/all/:uuid/:channel", webrenderer.HTTPAPIPlayAll)
		public.StaticFS("/static", http.Dir(datamodel.Storage.ServerHTTPDir()+"/static"))
	}
	/*
		Stream Control elements
	*/

	private.GET("/streams", webstreamer.HTTPAPIServerStreams)
	private.POST("/stream/:uuid/add", webstreamer.HTTPAPIServerStreamAdd)
	private.POST("/stream/:uuid/edit", webstreamer.HTTPAPIServerStreamEdit)
	private.GET("/stream/:uuid/delete", webstreamer.HTTPAPIServerStreamDelete)
	private.GET("/stream/:uuid/reload", webstreamer.HTTPAPIServerStreamReload)
	private.GET("/stream/:uuid/info", webstreamer.HTTPAPIServerStreamInfo)

	/*
		Streams Multi Control elements
	*/

	private.POST("/streams/multi/control/add", webstreamer.HTTPAPIServerStreamsMultiControlAdd)
	private.POST("/streams/multi/control/delete", webstreamer.HTTPAPIServerStreamsMultiControlDelete)

	/*
		Stream Channel elements
	*/

	private.POST("/stream/:uuid/channel/:channel/add", webstreamer.HTTPAPIServerStreamChannelAdd)
	private.POST("/stream/:uuid/channel/:channel/edit", webstreamer.HTTPAPIServerStreamChannelEdit)
	private.GET("/stream/:uuid/channel/:channel/delete", webstreamer.HTTPAPIServerStreamChannelDelete)
	private.GET("/stream/:uuid/channel/:channel/codec", webstreamer.HTTPAPIServerStreamChannelCodec)
	private.GET("/stream/:uuid/channel/:channel/reload", webstreamer.HTTPAPIServerStreamChannelReload)
	private.GET("/stream/:uuid/channel/:channel/info", webstreamer.HTTPAPIServerStreamChannelInfo)

	/*
		Stream video elements
	*/
	// HLS
	public.GET("/stream/:uuid/channel/:channel/hls/live/index.m3u8", mediastreamer.HTTPAPIServerStreamHLSM3U8)
	public.GET("/stream/:uuid/channel/:channel/hls/live/segment/:seq/file.ts", mediastreamer.HTTPAPIServerStreamHLSTS)
	// HLS LL
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/index.m3u8", mediastreamer.HTTPAPIServerStreamHLSLLM3U8)
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/init.mp4", mediastreamer.HTTPAPIServerStreamHLSLLInit)
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/segment/:segment/:any", mediastreamer.HTTPAPIServerStreamHLSLLM4Segment)
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/fragment/:segment/:fragment/:any", mediastreamer.HTTPAPIServerStreamHLSLLM4Fragment)
	// MSE
	public.GET("/stream/:uuid/channel/:channel/mse", mediastreamer.HTTPAPIServerStreamMSE)
	public.POST("/stream/:uuid/channel/:channel/webrtc", mediastreamer.HTTPAPIServerStreamWebRTC)

	/*
		HTTPS Mode Cert
		# Key considerations for algorithm "RSA" ≥ 2048-bit
		openssl genrsa -out server.key 2048

		# Key considerations for algorithm "ECDSA" ≥ secp384r1
		# List ECDSA the supported curves (openssl ecparam -list_curves)
		#openssl ecparam -genkey -name secp384r1 -out server.key
		#Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)

		openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	*/
	if datamodel.Storage.ServerHTTPS() {
		if datamodel.Storage.ServerHTTPSAutoTLSEnable() {
			go func() {
				err := autotls.Run(public, datamodel.Storage.ServerHTTPSAutoTLSName()+datamodel.Storage.ServerHTTPSPort())
				if err != nil {
					utils.Logger.Errorln("Start HTTPS Server Error", err)
				}
			}()
		} else {
			go func() {
				err := public.RunTLS(datamodel.Storage.ServerHTTPSPort(), datamodel.Storage.ServerHTTPSCert(), datamodel.Storage.ServerHTTPSKey())
				if err != nil {
					utils.Logger.Fatalln(err.Error())
					os.Exit(1)
				}
			}()
		}
	}

	err := public.Run(datamodel.Storage.ServerHTTPPort())
	if err != nil {
		utils.Logger.Fatal(err.Error())
		os.Exit(1)
	}
}

// CrossOrigin Access-Control-Allow-Origin any methods
func crossOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
