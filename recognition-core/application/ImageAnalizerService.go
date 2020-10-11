package recognitionapplication

import (
	"fmt"
	"go-intelligent-monitoring-system/domain"
	recognitionapplicationportout "go-intelligent-monitoring-system/recognition-core/application/portout"
	"os"

	log "github.com/sirupsen/logrus"
)

//ImageAnalizerService send images to recognition port to be analized.
type ImageAnalizerService struct {
	analizeAdapter      recognitionapplicationportout.ImageRecognitionPort
	notificationAdapter recognitionapplicationportout.NotificationPort
}

//AnalizeImage analize if faces on image are recognized or not
func (ias *ImageAnalizerService) AnalizeImage(image domain.Image) error {

	analizedImage, err := ias.analizeAdapter.Recognize(image)
	if err != nil {
		return err
	}

	if analizedImage.PersonNameDetected == "" {
		notification := createNotification(image)
		ias.notificationAdapter.NotifyTopic(notification)
		log.WithFields(log.Fields{"notification.Message": notification.Message, "notification.Topic": notification.Topic, "notification.Type": notification.Type, "analizedImage.Name": analizedImage.Name, "analizedImage.RecognitionCoreResponse": string(analizedImage.RecognitionCoreResponse)}).Info("Image correctly analized but Person is not Authorized")
	} else {
		log.WithFields(log.Fields{"analizedImage.PersonNameDetected": analizedImage.PersonNameDetected, "analizedImage.Name": analizedImage.Name, "analizedImage.RecognitionCoreResponse": string(analizedImage.RecognitionCoreResponse)}).Info("Image correctly analized and Person is Authorized")
	}

	return nil
}

//NewImageAnalizerService ...
func NewImageAnalizerService(analizeAdapter recognitionapplicationportout.ImageRecognitionPort, notificationAdapter recognitionapplicationportout.NotificationPort) *ImageAnalizerService {

	ias := &ImageAnalizerService{
		analizeAdapter:      analizeAdapter,
		notificationAdapter: notificationAdapter,
	}

	return ias
}

func createNotification(image domain.Image) domain.Notification {
	var notification domain.Notification
	notification.Image = image

	notification.Topic = os.Getenv("SNS_TOPIC")
	notification.Type = "AWS_SNS_TOPIC"

	//TODO: add S3 image Url
	notification.Message = fmt.Sprintf("The person detected is not in your people authorize collection. Image analized name: %s", image.Name)

	return notification
}
