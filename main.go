package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	openapi "github.com/vtpl1/go_rtsp_to_web/openapi/go"
	"github.com/vtpl1/go_rtsp_to_web/utils"
)

func main() {
	utils.InitializeLogger()
	utils.Logger.Info("Server CORE start")
	YojakaApiService := openapi.NewYojakaApiService()
	YojakaApiController := openapi.NewYojakaApiController(YojakaApiService)

	router := openapi.NewRouter(YojakaApiController)
	Storage    = NewStreamCore()
	go Storage.StreamChannelRunAll()
	go http.ListenAndServe(":8080", router)

	signalChanel := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChanel
		utils.Logger.Info("Server received signal", sig)
		done <- true
	}()
	utils.Logger.Info("Server started and waiting for signals")
	<-done
	time.Sleep(2 * time.Second)
	utils.Logger.Info("Server stopped working due to signal")
}

