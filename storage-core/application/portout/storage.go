package storageapplicationportout

import "go-intelligent-monitoring-system/domain"

//StorageImagePort ...
type StorageImagePort interface {
	Save(image domain.Image) error
}
