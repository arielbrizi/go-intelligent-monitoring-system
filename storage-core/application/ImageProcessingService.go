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
func (ips *ImageProcessingService) ProcessImage(imgData []byte) (string, error) {
	//TODO discard if image has not faces

	image := &domain.Image{
		Bytes: imgData,
		//TODO
	}

	ips.image2QueueService.SendImage2Queue(*image)

	ips.storageImageAdapter.Save(*image)

	return "", nil
}

//NewImageProcessingService ...
func NewImageProcessingService() *ImageProcessingService {
	return &ImageProcessingService{}
}
