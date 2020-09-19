package storageapplication

import (
	"go-intelligent-monitoring-system/domain"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
)

//Image2QueueService send the image to the queue to be processed
type Image2QueueService struct {
	imageToQueueAdapter storageapplicationportout.QueueImagePort
}

//SendImage2Queue ...
func (i2q *Image2QueueService) SendImage2Queue(image domain.Image) error {
	i2q.imageToQueueAdapter.SendImage2Queue(image)
	return nil
}

//NewImage2QueueService ...
func NewImage2QueueService(imageToQueueAdapter storageapplicationportout.QueueImagePort) *Image2QueueService {
	return &Image2QueueService{
		imageToQueueAdapter: imageToQueueAdapter,
	}
}
