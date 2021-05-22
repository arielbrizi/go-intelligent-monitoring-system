package configurationapplication

import (
	configurationadapterout "go-intelligent-monitoring-system/configuration-core/adapterout"
	configurationapplicationportin "go-intelligent-monitoring-system/configuration-core/application/portin"
	configurationapplicationportout "go-intelligent-monitoring-system/configuration-core/application/portout"
	storageadapterout "go-intelligent-monitoring-system/storage-core/adapterout"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"

	"io/ioutil"

	"testing"
)

func TestAddAuthorizedFace(t *testing.T) {

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = configurationadapterout.NewRekoAdapterTest()

	//Define the "Adapter Out" to be used to connect to the storage core
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = storageadapterout.NewImage2S3AdapterTest()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = NewFaceIndexerService(storageImageAdapter, rekoAdapter)

	fileBytes, err := ioutil.ReadFile("../../test/images/withFaces/3.jpg")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	authorizedFace, err := faceIndexerService.AddAuthorizedFace(fileBytes, "3.jpg", "", "")
	if err != nil {
		t.Errorf("Error adding authorized face file: %v", err)
	}

	if authorizedFace.Name != "3.jpg" {
		t.Errorf("Error on authorizedFace.Name. Should be 3.jpg")
	}

	if len(authorizedFace.Bytes) != len(fileBytes) {
		t.Errorf("Error on authorizedFace.Bytes. Bytes are different than fileBytes")
	}
}
