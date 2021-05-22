package recognitionapplication

import (
	"go-intelligent-monitoring-system/domain"
	recognitionadapterout "go-intelligent-monitoring-system/recognition-core/adapterout"
	recognitionapplicationportin "go-intelligent-monitoring-system/recognition-core/application/portin"
	recognitionapplicationportout "go-intelligent-monitoring-system/recognition-core/application/portout"
	"io/ioutil"
	"testing"
)

func TestAnalizeImage(t *testing.T) {
	//Define the "Adapter Out" to be used to connect to the recognition core
	var analizeAdapter recognitionapplicationportout.ImageRecognitionPort
	analizeAdapter = recognitionadapterout.NewRekoAdapterTest()

	//Define the "Adapter Out" to be used to connect to notification core
	var notificationAdapter recognitionapplicationportout.NotificationPort
	notificationAdapter = recognitionadapterout.NewSNSAdapterTest()

	//Define the "Adapter Out" to be used to save categorized images (authorized, not authorized, etc)
	var imageStorageAdapter recognitionapplicationportout.ImageStoragePort
	imageStorageAdapter = recognitionadapterout.NewFtpImageStorageAdapterTest()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var imageAnalizerService recognitionapplicationportin.QueueImagePort
	imageAnalizerService = NewImageAnalizerService(analizeAdapter, notificationAdapter, imageStorageAdapter)

	fileBytes, _ := ioutil.ReadFile("../../test/images/withFaces/3.jpg")

	var image domain.Image
	image.Name = "3.jpg"
	image.Bytes = fileBytes

	analizedImage, errAnalize := imageAnalizerService.AnalizeImage(image)
	if errAnalize != nil {
		t.Errorf("Error Analizing Image %v", errAnalize)
	}
	if analizedImage.PersonNameDetected == "" {
		t.Errorf("Person not detected and is a known face.")
	}

}
