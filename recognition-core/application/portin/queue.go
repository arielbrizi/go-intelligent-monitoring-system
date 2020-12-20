package recognitionapplicationportin

import "go-intelligent-monitoring-system/domain"

//QueueImagePort ...
type QueueImagePort interface {
	AnalizeImage(image domain.Image) (*domain.AnalizedImage, error)
}
