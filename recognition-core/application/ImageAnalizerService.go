package recognitionapplication

import (
	"go-intelligent-monitoring-system/domain"
)

//ImageAnalizerService send images to recognition port to be analized.
type ImageAnalizerService struct {
}

//AnalizeImage analize if faces on image are recognized or not
func (ias *ImageAnalizerService) AnalizeImage(imgData domain.Image) error {

	//TODO: analize image
	return nil
}

//NewImageAnalizerService ...
func NewImageAnalizerService() *ImageAnalizerService {

	ias := &ImageAnalizerService{}

	return ias
}
