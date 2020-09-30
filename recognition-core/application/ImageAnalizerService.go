package recognitionapplication

import (
	"fmt"
	"go-intelligent-monitoring-system/domain"
	recognitionapplicationportout "go-intelligent-monitoring-system/recognition-core/application/portout"
)

//ImageAnalizerService send images to recognition port to be analized.
type ImageAnalizerService struct {
	analizeAdapter recognitionapplicationportout.ImageRecognitionPort
}

//AnalizeImage analize if faces on image are recognized or not
func (ias *ImageAnalizerService) AnalizeImage(image domain.Image) error {

	analizedImage, err := ias.analizeAdapter.Recognize(image)
	if err != nil {
		return err
	}

	fmt.Printf("Image '%s' analized. Name of Person Detected: '%s'", analizedImage.Name, analizedImage.PersonNameDetected)

	//TODO: process alarm port

	return nil
}

//NewImageAnalizerService ...
func NewImageAnalizerService(analizeAdapter recognitionapplicationportout.ImageRecognitionPort) *ImageAnalizerService {

	ias := &ImageAnalizerService{
		analizeAdapter: analizeAdapter,
	}

	return ias
}
