package main

import (
	"fmt"
	storageadapterin "go-intelligent-monitoring-system/storage-core/adapterin"
	storageapplicationportin "go-intelligent-monitoring-system/storage-core/application/portin"
)

func main() {
	fmt.Println("Initializing Intelligent Monitoring System")

	var filterImageService storageapplicationportin.InputPort
	//TODO Instanciar filterImageService

	ftpToInputAdapter := storageadapterin.NewFtpToInputAdapter(filterImageService)

	ftpToInputAdapter.ProcessImages()
}
