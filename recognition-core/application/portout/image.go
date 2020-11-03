package recognitionapplicationportout

import "go-intelligent-monitoring-system/domain"

//ImageStoragePort for save categorized images
type ImageStoragePort interface {
	SaveAuthorizedImage(image domain.Image) error
	SaveNotAuthorizedImage(image domain.Image) error
}
