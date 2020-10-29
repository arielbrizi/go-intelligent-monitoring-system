package main

import (
	"go-intelligent-monitoring-system/utils"
	"io"
	"os"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel) // TODO: get from Config

	var file, _ = os.OpenFile("logFile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(io.MultiWriter(file, os.Stdout))
	log.AddHook(logrus_stack.StandardHook()) //Sets the time and code line on each log.

	log.Info("Initializing Intelligent Monitoring System")

	// --------------- Configuration Port IN --------------- //

	confDirectoryAdapter := utils.NewConfDirectoryAdapter()
	confDirectoryAdapter.AddAuthorizedFaces()

	// --------------- Input Port IN--------------- //

	ftpToInputAdapter := utils.NewftpToInputAdapter()
	go ftpToInputAdapter.Process()

	// --------------- Queue PORT IN --------------- //

	queueInAdapter := utils.NewQueueInAdapter()
	queueInAdapter.ReceiveImagesFromQueue()

}
