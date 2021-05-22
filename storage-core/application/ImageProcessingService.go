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
	Image2QueueService  Image2QueueService
	StorageImageAdapter storageapplicationportout.StorageImagePort
}

//ProcessImage analize if image has faces to send it to the Queue.
func (ips *ImageProcessingService) ProcessImage(imgData []byte, fileName string) error {

	faces, err := utils.FacesOnImage(imgData)
	if err != nil {
		return err
	}

	if faces == 0 {
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

	return ips
}
