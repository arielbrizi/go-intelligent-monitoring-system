package recognitionapplicationportout

import "go-intelligent-monitoring-system/domain"

//ImageRecognitionPort ...
type ImageRecognitionPort interface {
	Recognize(image domain.Image) (*domain.AnalizedImage, error)
}
