package main

import (
	"fmt"
	storageadapterin "go-intelligent-monitoring-system/storage-core/adapterin"
	storageapplication "go-intelligent-monitoring-system/storage-core/application"
	storageapplicationportin "go-intelligent-monitoring-system/storage-core/application/portin"
)

func main() {
	fmt.Println("Initializing Intelligent Monitoring System")

	var imageProcessingService storageapplicationportin.InputImagePort
	var video2ImageService storageapplicationportin.InputVideoPort

	imageProcessingService = storageapplication.NewImageProcessingService()
	video2ImageService = storageapplication.NewVideo2ImageService()

	ftpToInputAdapter := storageadapterin.NewFtpToInputAdapter(imageProcessingService, video2ImageService)

	ftpToInputAdapter.Process()
}
