package main

import (
	"go-intelligent-monitoring-system/utils"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel) // TODO: get from Config
	log.AddHook(logrus_stack.StandardHook())

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
