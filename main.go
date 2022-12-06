package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var logger *logrus.Entry

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
		DisableColors: false,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	logger := log.WithFields(logrus.Fields{
		"module": "main",
		"func":   "main",
	})
	logger.Info("Server CORE start")

	go Storage.StreamChannelRunAll()

	signalChanel := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChanel
		logger.Info("Server received signal", sig)
		done <- true
	}()
	logger.Info("Server started and waiting for signals")
	<-done
	time.Sleep(2 * time.Second)
	logger.Info("Server stopped working due to signal")
}
