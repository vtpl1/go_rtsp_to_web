package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vtpl1/go_rtsp_to_web/utils"
)

func main() {
	utils.InitializeLogger()
	utils.Logger.Info("Server CORE start")

	go Storage.StreamChannelRunAll()

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
