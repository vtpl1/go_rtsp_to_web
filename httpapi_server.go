package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/vtpl1/go_rtsp_to_web/utils"
)

func HTTPAPIServer() {
	utils.Logger.Info("HTTP Server start")
	var public *gin.Engine
	if !Storage.ServerHTTPDebug() {
		gin.SetMode(gin.ReleaseMode)
		public = gin.New()
	} else {
		gin.SetMode(gin.DebugMode)
		public = gin.Default()
	}
	public.Use(crossOrigin())
	// Add private login password protect methods
	private := public.Group("/")
	if Storage.ServerHTTPLogin() != "" && Storage.ServerHTTPPassword() != "" {
		private.Use(gin.BasicAuth(gin.Accounts{Storage.ServerHTTPLogin(): Storage.ServerHTTPPassword()}))
	}

		/*
		Static HTML Files Demo Mode
	*/

	if Storage.ServerHTTPDemo() {
		public.LoadHTMLGlob(Storage.ServerHTTPDir() + "/templates/*")
		public.GET("/", HTTPAPIServerIndex)
		public.GET("/pages/stream/list", HTTPAPIStreamList)
		public.GET("/pages/stream/add", HTTPAPIAddStream)
		public.GET("/pages/stream/edit/:uuid", HTTPAPIEditStream)
		public.GET("/pages/player/hls/:uuid/:channel", HTTPAPIPlayHls)
		public.GET("/pages/player/mse/:uuid/:channel", HTTPAPIPlayMse)
		public.GET("/pages/player/webrtc/:uuid/:channel", HTTPAPIPlayWebrtc)
		public.GET("/pages/multiview", HTTPAPIMultiview)
		public.Any("/pages/multiview/full", HTTPAPIFullScreenMultiView)
		public.GET("/pages/documentation", HTTPAPIServerDocumentation)
		public.GET("/pages/player/all/:uuid/:channel", HTTPAPIPlayAll)
		public.StaticFS("/static", http.Dir(Storage.ServerHTTPDir()+"/static"))
	}

	
	err := public.Run(Storage.ServerHTTPPort())
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
