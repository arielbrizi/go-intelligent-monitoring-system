package storageapplication

import (
	"errors"
	"go-intelligent-monitoring-system/domain"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
	"go-intelligent-monitoring-system/storage-core/application/utils"
	"os"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

//ImageProcessingService filter image without faces on it and send it to the Queue
type ImageProcessingService struct {
	Image2QueueService  Image2QueueService
	StorageImageAdapter storageapplicationportout.StorageImagePort
	redisClient         *redis.Client
}

//ProcessImage analize if image has faces to send it to the Queue.
func (ips *ImageProcessingService) ProcessImage(imgData []byte, fileName string) error {

	faces, err := utils.FacesOnImage(imgData)
	if err != nil {
		return err
	}

	if faces == 0 {
		log.WithFields(log.Fields{"fileName": fileName}).Info("No faces on image")
		deleteImage(fileName)
		return nil
	}

	statusRecognition, err := ips.redisClient.Get("statusRecognition").Result()
	if err != nil {
		log.WithError(err).Error("error getting redis value: statusRecognition")
		return nil
	}
	if statusRecognition == "OFF" {
		log.WithFields(log.Fields{"fileName": fileName}).Info("statusRecognition == OFF")
		return nil
	}

	var bucket string
	if bucket = os.Getenv("CAMARA_DOMAIN"); bucket == "" {
		return errors.New("CAMARA_DOMAIN env not defined")
	}

	if imgData == nil || len(imgData) == 0 {
		return errors.New("ProcessImage: image '" + fileName + "' empty")
	}

	image := &domain.Image{
		Bytes:  imgData,
		Name:   fileName,
		Bucket: bucket,
		//TODO: complete all image attributes
	}

	errSave := ips.StorageImageAdapter.Save(*image)
	if errSave != nil {
		return errSave
	}

	urlImage, _ := ips.StorageImageAdapter.GetURL(*image)
	image.URL = urlImage

	errSend := ips.Image2QueueService.SendImage2Queue(*image)
	if errSend != nil {
		return errSend
	}

	return nil
}

//NewImageProcessingService ...
func NewImageProcessingService(storageImageAdapter storageapplicationportout.StorageImagePort, queueAdapter storageapplicationportout.QueueImagePort) *ImageProcessingService {
	i2qs := NewImage2QueueService(queueAdapter)

	ips := &ImageProcessingService{
		StorageImageAdapter: storageImageAdapter,
		Image2QueueService:  *i2qs,
	}

	ips.redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis-server:6379",
		Password: os.Getenv("REDIS_PASS"),
		DB:       0, // use default DB
	})

	return ips
}

func deleteImage(name string) {
	var t = time.Now()
	var today = t.Format("20060102")
	var ftpDirectory = os.Getenv("FTP_DIRECTORY")
	var ftpTodayDirectory = ftpDirectory + today + "/"
	err := os.Remove(ftpTodayDirectory + name)
	if err != nil {
		log.WithError(err).Error("error deleting file")
		return
	}
	log.WithFields(log.Fields{"fileName": name}).Info("deleted file:" + ftpTodayDirectory + name)

}
