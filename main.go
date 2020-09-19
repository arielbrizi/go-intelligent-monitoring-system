package main

import (
	"fmt"
	storageadapterin "go-intelligent-monitoring-system/storage-core/adapterin"
	storageadapterout "go-intelligent-monitoring-system/storage-core/adapterout"
	storageapplication "go-intelligent-monitoring-system/storage-core/application"
	storageapplicationportin "go-intelligent-monitoring-system/storage-core/application/portin"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
)

func main() {
	fmt.Println("Initializing Intelligent Monitoring System")

	var imageProcessingService storageapplicationportin.InputImagePort
	var video2ImageService storageapplicationportin.InputVideoPort

	//Storage settings: S3 Adapter
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = storageadapterout.NewImage2S3Adapter()

	//Queue settings: Kafka Adapter
	var queueAdapter storageapplicationportout.QueueImagePort
	queueAdapter = storageadapterout.NewKafkaAdapter()

	imageProcessingService = storageapplication.NewImageProcessingService(storageImageAdapter, queueAdapter)
	video2ImageService = storageapplication.NewVideo2ImageService()

	//Input settings: Ftp Adapter
	ftpToInputAdapter := storageadapterin.NewFtpToInputAdapter(imageProcessingService, video2ImageService)

	ftpToInputAdapter.Process()
}
