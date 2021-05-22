package storageapplicationportout

import "go-intelligent-monitoring-system/domain"

//StorageImagePort ...
type StorageImagePort interface {
	Save(image domain.Image) error
	GetURL(image domain.Image) (string, error)
}
