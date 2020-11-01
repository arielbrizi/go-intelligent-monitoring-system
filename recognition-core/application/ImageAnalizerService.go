package recognitionapplication

import (
	"fmt"
	"go-intelligent-monitoring-system/domain"
	recognitionapplicationportout "go-intelligent-monitoring-system/recognition-core/application/portout"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//ImageAnalizerService send images to recognition port to be analized.
type ImageAnalizerService struct {
	analizeAdapter      recognitionapplicationportout.ImageRecognitionPort
	notificationAdapter recognitionapplicationportout.NotificationPort
	snsTopic            string
	ftpDirectory        string
}

//AnalizeImage analize if faces on image are recognized or not
func (ias *ImageAnalizerService) AnalizeImage(image domain.Image) error {

	analizedImage, err := ias.analizeAdapter.Recognize(image)
	if err != nil {
		return err
	}

	if analizedImage.PersonNameDetected == "" {
		notification := ias.createNotification(image)
		ias.notificationAdapter.NotifyTopic(notification)
		log.WithFields(log.Fields{"notification.Message": notification.Message, "notification.Topic": notification.Topic, "notification.Type": notification.Type, "analizedImage.Name": analizedImage.Name, "analizedImage.RecognitionCoreResponse": string(analizedImage.RecognitionCoreResponse)}).Info("Image correctly analized but Person is not Authorized")
		ias.moveFaceNotRecognized(image)
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
		snsTopic:            os.Getenv("SNS_TOPIC"),
		ftpDirectory:        os.Getenv("FTP_DIRECTORY"),
	}

	return ias
}

func (ias *ImageAnalizerService) createNotification(image domain.Image) domain.Notification {
	var notification domain.Notification
	notification.Image = image

	notification.Topic = ias.snsTopic
	notification.Type = "AWS_SNS_TOPIC"

	//TODO: add S3 image Url
	notification.Message = fmt.Sprintf("The person detected is not in your people authorize collection. Image analized name: %s", image.Name)

	return notification
}

//moveFaceNotRecognized move image to faces not authorized directory
func (ias *ImageAnalizerService) moveFaceNotRecognized(image domain.Image) {

	var t = time.Now()
	var today = t.Format("20060102")

	var ftpTodayDirectory = ias.ftpDirectory + today + "/"
	var ftpTodayDirectoryProcessed = strings.Replace(ftpTodayDirectory, today, today+"_processed", 1)
	var ftpTodayDirectoryFacesNotAuth = strings.Replace(ftpTodayDirectory, today, today+"_faces_not_auth", 1)

	err := os.Rename(ftpTodayDirectoryProcessed+image.Name, ftpTodayDirectoryFacesNotAuth+image.Name)
	if err != nil {
		log.WithFields(log.Fields{"ftpTodayDirectoryProcessed": ftpTodayDirectoryProcessed, "ftpTodayDirectoryFacesNotAuth": ftpTodayDirectoryFacesNotAuth, "fileName": image.Name}).WithError(err).Error("Error moving file to Not Authorized directory")
	}

}
