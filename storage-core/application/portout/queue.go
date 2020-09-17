package storageapplicationportout

import "go-intelligent-monitoring-system/domain"

//QueueImagePort...
type QueueImagePort interface {
	SendImage2Queue(image domain.Image) (string, error)
}
