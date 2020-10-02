package main

import (
	"fmt"
	recognitionadapterin "go-intelligent-monitoring-system/recognition-core/adapterin"
	recognitionadapterout "go-intelligent-monitoring-system/recognition-core/adapterout"
	recognitionapplication "go-intelligent-monitoring-system/recognition-core/application"
	recognitionapplicationportin "go-intelligent-monitoring-system/recognition-core/application/portin"
	recognitionapplicationportout "go-intelligent-monitoring-system/recognition-core/application/portout"
	storageadapterin "go-intelligent-monitoring-system/storage-core/adapterin"
	storageadapterout "go-intelligent-monitoring-system/storage-core/adapterout"
	storageapplication "go-intelligent-monitoring-system/storage-core/application"
	storageapplicationportin "go-intelligent-monitoring-system/storage-core/application/portin"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
)

func main() {
	fmt.Println("Initializing Intelligent Monitoring System")

	// --------------- Input Port --------------- //

	//IMPLEMENTATION --> Storage settings: S3 Adapter
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = storageadapterout.NewImage2S3Adapter()

	//IMPLEMENTATION --> Queue settings: Kafka Adapter Out
	var queueAdapterOut storageapplicationportout.QueueImagePort
	queueAdapterOut = storageadapterout.NewKafkaAdapter()

	var imageProcessingService storageapplicationportin.InputImagePort
	imageProcessingService = storageapplication.NewImageProcessingService(storageImageAdapter, queueAdapterOut)

	var video2ImageService storageapplicationportin.InputVideoPort
	video2ImageService = storageapplication.NewVideo2ImageService()

	//Input settings: Ftp Adapter
	ftpToInputAdapter := storageadapterin.NewFtpToInputAdapter(imageProcessingService, video2ImageService)

	go ftpToInputAdapter.Process()

	// --------------- QUEUE PORT IN --------------- //

	var analizeAdapter recognitionapplicationportout.ImageRecognitionPort
	analizeAdapter = recognitionadapterout.NewRekoAdapter()

	var notificationAdapter recognitionapplicationportout.NotificationPort
	notificationAdapter = recognitionadapterout.NewSNSAdapter()

	var imageAnalizerService recognitionapplicationportin.QueueImagePort
	imageAnalizerService = recognitionapplication.NewImageAnalizerService(analizeAdapter, notificationAdapter)

	//Input Queue settings: Kafka Adapter In
	queueInAdapter := recognitionadapterin.NewKafkaAdapter(imageAnalizerService)

	queueInAdapter.ReceiveImagesFromQueue()
}
