package main

import (
	"go-intelligent-monitoring-system/utils"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

// @title Intelligent Monitoring System
// @version 1.0
// @description Configuration and Analize Services.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://arielbrizi.github.io/go-intelligent-monitoring-system/
// @contact.email arielbrizi@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel) // TODO: get from Config

	//One log file per day. Purge logs older than 7 days
	rl, _ := rotatelogs.New("log.%Y%m%d", rotatelogs.WithMaxAge(-1), rotatelogs.WithRotationCount(7))
	log.SetOutput(rl)

	log.AddHook(logrus_stack.StandardHook()) //Sets the time and code line on each log.
	log.Info("Initializing Intelligent Monitoring System")

	// --------------- Configuration Core - Port IN --------------- //

	confDirectoryAdapter := utils.NewConfDirectoryAdapter()
	confDirectoryAdapter.AddAuthorizedFaces()

	// --------------- Storage Core - Port IN --------------- //

	ftpToInputAdapter := utils.NewftpToInputAdapter()
	go ftpToInputAdapter.Process()

	// --------------- Recognition Core - PORT IN --------------- //

	queueInAdapter := utils.NewQueueInAdapter()
	go queueInAdapter.ReceiveImagesFromQueue()

	// --------------- API from All Cores ---------------- //

	utils.RunAPI()

}
