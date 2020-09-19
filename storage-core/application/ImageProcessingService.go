package storageapplication

import (
	"go-intelligent-monitoring-system/domain"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
)

//ImageProcessingService filter image without faces on it and send it to the Queue
type ImageProcessingService struct {
	image2QueueService  Image2QueueService
	storageImageAdapter storageapplicationportout.StorageImagePort
}

//ProcessImage analize if image has faces to send it to the Queue.
func (ips *ImageProcessingService) ProcessImage(imgData []byte) error {
	//TODO discard if image has not faces

	image := &domain.Image{
		Bytes: imgData,
		//TODO
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
