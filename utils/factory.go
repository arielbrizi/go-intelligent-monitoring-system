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
	"os"
	"sync"

	_ "go-intelligent-monitoring-system/docs" // docs generated by Swag CLI

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var rekoAdapter *configurationadapterout.RekoAdapter
var image2S3Adapter *storageadapterout.Image2S3Adapter
var faceIndexerService *configurationapplication.FaceIndexerService

var once sync.Once
var once2 sync.Once
var once3 sync.Once

//NewConfDirectoryAdapter ...
func NewConfDirectoryAdapter() *configurationadapterin.DirectoryAdapter {

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = GetRekoAdapter()

	//Define the "Adapter Out" to be used to connect to the storage core
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = GetImage2S3Adapter()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = GetFaceIndexerService(storageImageAdapter, rekoAdapter)

	//"Adapter In": DirectoryAdapter gets the authorized faces from a directory
	confDirectoryAdapter := configurationadapterin.NewDirectoryAdapter(faceIndexerService)

	return confDirectoryAdapter

}

//NewftpToInputAdapter ...
func NewftpToInputAdapter() *storageadapterin.FtpToInputAdapter {

	//Define the "Adapter Out" to be used to connect to the storage core
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = GetImage2S3Adapter()

	//Define the "Adapter Out" to be used to connect to the Image Queue core: Kafka
	var queueAdapterOut storageapplicationportout.QueueImagePort
	queueAdapterOut = storageadapterout.NewKafkaAdapter()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out" for Images
	var imageProcessingService storageapplicationportin.InputImagePort
	imageProcessingService = storageapplication.NewImageProcessingService(storageImageAdapter, queueAdapterOut)

	//"Adapter In": FtpToInputAdapter gets the images to be analized from FTP directory
	return storageadapterin.NewFtpToInputAdapter(imageProcessingService)

}

//NewQueueInAdapter ...
func NewQueueInAdapter() *recognitionadapterin.KafkaAdapter {

	//Define the "Adapter Out" to be used to connect to the recognition core
	var analizeAdapter recognitionapplicationportout.ImageRecognitionPort
	analizeAdapter = recognitionadapterout.NewRekoAdapter()

	//Define the "Adapter Out" to be used to connect to notification core
	var notificationAdapter recognitionapplicationportout.NotificationPort
	switch os.Getenv("NOTIFICATION_TYPE") {
	case "TELEGRAM":
		//The Telegram Adapter have input (telegram commands) and output methods (send messages)
		notificationAdapter = recognitionadapterin.NewTelegramAdapter()
	case "SNS":
		notificationAdapter = recognitionadapterout.NewSNSAdapter()

	}

	//Define the "Adapter Out" to be used to save categorized images (authorized, not authorized, etc)
	var imageStorageAdapter recognitionapplicationportout.ImageStoragePort
	imageStorageAdapter = recognitionadapterout.NewFtpImageStorageAdapter()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var imageAnalizerService recognitionapplicationportin.QueueImagePort
	imageAnalizerService = recognitionapplication.NewImageAnalizerService(analizeAdapter, notificationAdapter, imageStorageAdapter)

	//"Adapter In": KafkaAdapter gets the images to be analized from Kafka
	return recognitionadapterin.NewKafkaAdapter(imageAnalizerService)
}

//NewConfAPIAdapter ...
func NewConfAPIAdapter() *configurationadapterin.APIAdapter {

	//Define the "Adapter Out" to be used to connect to the storage core
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = GetImage2S3Adapter()

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = GetRekoAdapter()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = GetFaceIndexerService(storageImageAdapter, rekoAdapter)

	//"Adapter In": APIAdapter sets the authorized face from the request
	confAPIAdapter := configurationadapterin.NewAPIAdapter(faceIndexerService)

	return confAPIAdapter

}

func addConfigurationCoreAPI(r *gin.Engine) {

	confAPIAdapter := NewConfAPIAdapter()

	r.GET("/configuration-core/authorized-face/:collectionName", confAPIAdapter.GetAuthorizedFacesHandler)
	r.POST("/configuration-core/authorized-face", confAPIAdapter.AddAuthorizedFaceHandler)
	r.DELETE("/configuration-core/authorized-face", confAPIAdapter.DeleteAuthorizedFaceHandler)
}

//Return RekoAdapter Singleton
func GetRekoAdapter() *configurationadapterout.RekoAdapter {

	once.Do(func() { // <-- atomic, does not allow repeating

		rekoAdapter = configurationadapterout.NewRekoAdapter()

	})

	return rekoAdapter

}

//Return Image2S3Adapter Singleton
func GetImage2S3Adapter() *storageadapterout.Image2S3Adapter {

	once2.Do(func() { // <-- atomic, does not allow repeating

		image2S3Adapter = storageadapterout.NewImage2S3Adapter()
	})

	return image2S3Adapter

}

//Return FaceIndexerService Singleton
func GetFaceIndexerService(storageImageAdapter storageapplicationportout.StorageImagePort, rekoAdapter configurationapplicationportout.ImageRecognitionPort) *configurationapplication.FaceIndexerService {

	once3.Do(func() { // <-- atomic, does not allow repeating

		faceIndexerService = configurationapplication.NewFaceIndexerService(storageImageAdapter, rekoAdapter)
	})

	return faceIndexerService

}

//RunAPI ...
func RunAPI() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	addConfigurationCoreAPI(r)

	//To publish swagger doc (Installation instructions on https://github.com/swaggo/gin-swagger)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
