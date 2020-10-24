package storageapplication

import (
	"errors"
	"go-intelligent-monitoring-system/domain"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
	"go-intelligent-monitoring-system/storage-core/application/utils"
	"os"

	log "github.com/sirupsen/logrus"
)

//ImageProcessingService filter image without faces on it and send it to the Queue
type ImageProcessingService struct {
	image2QueueService  Image2QueueService
	storageImageAdapter storageapplicationportout.StorageImagePort
}

//ProcessImage analize if image has faces to send it to the Queue.
func (ips *ImageProcessingService) ProcessImage(imgData []byte, fileName string) error {

	if utils.FacesOnImage(imgData) == 0 {
		log.WithFields(log.Fields{"fileName": fileName}).Info("No faces on image")
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

	err := ips.storageImageAdapter.Save(*image)
	if err != nil {
		return err
	}

	err = ips.image2QueueService.SendImage2Queue(*image)
	if err != nil {
		return err
	}

	return nil
}

//NewImageProcessingService ...
func NewImageProcessingService(storageImageAdapter storageapplicationportout.StorageImagePort, queueAdapter storageapplicationportout.QueueImagePort) *ImageProcessingService {
	i2qs := NewImage2QueueService(queueAdapter)

	ips := &ImageProcessingService{
		storageImageAdapter: storageImageAdapter,
		image2QueueService:  *i2qs,
	}

	return ips
}
