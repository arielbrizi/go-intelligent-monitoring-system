package main

import (
	"fmt"
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

	logrus_stack "github.com/Gurpartap/logrus-stack"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel) // TODO: get from Config
	log.AddHook(logrus_stack.StandardHook())

	fmt.Println("Initializing Intelligent Monitoring System")

	// --------------- CONFIGURATION PORT IN --------------- //

	//IMPLEMENTATION --> Recognition settings: AWS Reko Adapter
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = configurationadapterout.NewRekoAdapter()

	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = configurationapplication.NewFaceIndexerService(rekoAdapter)

	//Input Configuration settings: Directory Adapter to add authorized faces saved on it (AUTHORIZED_FACES_DIRECTORY)
	confDirectoryAdapter := configurationadapterin.NewDirectoryAdapter(faceIndexerService)

	confDirectoryAdapter.AddAuthorizedFaces()

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

	//IMPLEMENTATION --> Recognition settings: AWS Reko Adapter
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
