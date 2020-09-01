package main

import (
	"fmt"
	storageadapterin "go-intelligent-monitoring-system/storage-core/adapterin"
	storageapplication "go-intelligent-monitoring-system/storage-core/application"
	storageapplicationportin "go-intelligent-monitoring-system/storage-core/application/portin"
)

func main() {
	fmt.Println("Initializing Intelligent Monitoring System")

	var filterImageService storageapplicationportin.InputPort
	filterImageService = storageapplication.NewFilterImageService()

	ftpToInputAdapter := storageadapterin.NewFtpToInputAdapter(filterImageService)

	ftpToInputAdapter.ProcessImages()
}
