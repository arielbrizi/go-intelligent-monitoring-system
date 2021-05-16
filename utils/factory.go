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

	_ "go-intelligent-monitoring-system/docs" // docs generated by Swag CLI

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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
	switch os.Getenv("NOTIFICATION_TYPE") {
	case "TELEGRAM":
		notificationAdapter = recognitionadapterout.NewTelegramAdapter()
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

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = configurationadapterout.NewRekoAdapter()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = configurationapplication.NewFaceIndexerService(rekoAdapter)

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