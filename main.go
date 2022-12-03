package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithFields(log.Fields{
		"module": "main",
		"func":   "main",
	}).Info("Server CORE start")
	log.Info("Hello, world.")
	signalChanel := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChanel
		log.WithFields(log.Fields{
			"module": "main",
			"func":   "main",
		}).Info("Server receive signal", sig)
		done <- true
	}()
	log.WithFields(log.Fields{
		"module": "main",
		"func":   "main",
	}).Info("Server start success a wait signals")
	<-done
	time.Sleep(2 * time.Second)
	log.WithFields(log.Fields{
		"module": "main",
		"func":   "main",
	}).Info("Server stop working by signal")
}
