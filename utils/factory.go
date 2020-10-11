package utils

import (
	configurationadapterin "go-intelligent-monitoring-system/configuration-core/adapterin"
	configurationadapterout "go-intelligent-monitoring-system/configuration-core/adapterout"
	configurationapplication "go-intelligent-monitoring-system/configuration-core/application"
	configurationapplicationportin "go-intelligent-monitoring-system/configuration-core/application/portin"
	configurationapplicationportout "go-intelligent-monitoring-system/configuration-core/application/portout"
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

//NewConfDirectoryAdapter ...
func NewConfDirectoryAdapter() *configurationadapterin.DirectoryAdapter {

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = configurationadapterout.NewRekoAdapter()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = configurationapplication.NewFaceIndexerService(rekoAdapter)

	//"Adapter In": DirectoryAdapter gets the authorized faces from a directory
	confDirectoryAdapter := configurationadapterin.NewDirectoryAdapter(faceIndexerService)

	return confDirectoryAdapter

}

//NewftpToInputAdapter ...
func NewftpToInputAdapter() *storageadapterin.FtpToInputAdapter {

	//Define the "Adapter Out" to be used to connect to the storage core
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = storageadapterout.NewImage2S3Adapter()

	//Define the "Adapter Out" to be used to connect to the Image Queue core: Kafka
	var queueAdapterOut storageapplicationportout.QueueImagePort
	queueAdapterOut = storageadapterout.NewKafkaAdapter()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out" for Images
	var imageProcessingService storageapplicationportin.InputImagePort
	imageProcessingService = storageapplication.NewImageProcessingService(storageImageAdapter, queueAdapterOut)

	//Define the service to be  used between the "Adapter In" and the "Adapter Out" for Videos
	var video2ImageService storageapplicationportin.InputVideoPort
	video2ImageService = storageapplication.NewVideo2ImageService()

	//"Adapter In": FtpToInputAdapter gets the images to be analized from FTP directory
	return storageadapterin.NewFtpToInputAdapter(imageProcessingService, video2ImageService)

}

//NewQueueInAdapter ...
func NewQueueInAdapter() *recognitionadapterin.KafkaAdapter {

	//Define the "Adapter Out" to be used to connect to the recognition core
	var analizeAdapter recognitionapplicationportout.ImageRecognitionPort
	analizeAdapter = recognitionadapterout.NewRekoAdapter()

	//Define the "Adapter Out" to be used to connect to notification core
	var notificationAdapter recognitionapplicationportout.NotificationPort
	notificationAdapter = recognitionadapterout.NewSNSAdapter()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var imageAnalizerService recognitionapplicationportin.QueueImagePort
	imageAnalizerService = recognitionapplication.NewImageAnalizerService(analizeAdapter, notificationAdapter)

	//"Adapter In": KafkaAdapter gets the images to be analized from Kafka
	return recognitionadapterin.NewKafkaAdapter(imageAnalizerService)
}
