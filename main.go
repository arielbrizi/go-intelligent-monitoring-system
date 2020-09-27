package main

import (
	"fmt"
	recognitionadapterin "go-intelligent-monitoring-system/recognition-core/adapterin"
	recognitionapplicationportin "go-intelligent-monitoring-system/recognition-core/application/portin"
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
	var queueAdapterOut storageapplicationportout.QueueImagePort
	queueAdapterOut = storageadapterout.NewKafkaAdapter()

	imageProcessingService = storageapplication.NewImageProcessingService(storageImageAdapter, queueAdapterOut)
	video2ImageService = storageapplication.NewVideo2ImageService()

	//Input settings: Ftp Adapter
	ftpToInputAdapter := storageadapterin.NewFtpToInputAdapter(imageProcessingService, video2ImageService)

	//Queue settings: Kafka Adapter
	var queueAdapterIn recognitionapplicationportin.QueueImagePort
	queueAdapterIn = recognitionadapterin.NewKafkaAdapter()

	go ftpToInputAdapter.Process()

	queueAdapterIn.ReceiveImagesFromQueue()
}
