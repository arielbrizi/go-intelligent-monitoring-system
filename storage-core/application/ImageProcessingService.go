package storageapplication

import (
	"errors"
	"go-intelligent-monitoring-system/domain"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
	"os"
)

//ImageProcessingService filter image without faces on it and send it to the Queue
type ImageProcessingService struct {
	image2QueueService  Image2QueueService
	storageImageAdapter storageapplicationportout.StorageImagePort
}

//ProcessImage analize if image has faces to send it to the Queue.
func (ips *ImageProcessingService) ProcessImage(imgData []byte, fileName string) error {
	//TODO discard if image has not faces

	var bucket string
	if bucket = os.Getenv("CAMARA_DOMAIN"); bucket == "" {
		return errors.New("CAMARA_DOMAIN env not defined")
	}

	if imgData == nil || len(imgData) == 0 {
		return errors.New("ProcessImage: image '" + fileName + "' empty")
	}

	image := &domain.Image{
		Bytes: imgData,
		Name: fileName,
		Bucket: bucket,
		//TODO: complete all image attributes
	}

	ips.storageImageAdapter.Save(*image)

	ips.image2QueueService.SendImage2Queue(*image)

	return nil
}

//NewImageProcessingService ...
func NewImageProcessingService(storageImageAdapter storageapplicationportout.StorageImagePort, queueAdapter storageapplicationportout.QueueImagePort) *ImageProcessingService {
	i2qs := NewImage2QueueService(queueAdapter)

	ips := &ImageProcessingService{
		storageImageAdapter: storageImageAdapter,
		image2QueueService: *i2qs,
	}

	return ips
}
