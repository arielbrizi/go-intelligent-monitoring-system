package configurationadapterin

import (
	configurationadapterout "go-intelligent-monitoring-system/configuration-core/adapterout"
	configurationapplication "go-intelligent-monitoring-system/configuration-core/application"
	configurationapplicationportin "go-intelligent-monitoring-system/configuration-core/application/portin"
	configurationapplicationportout "go-intelligent-monitoring-system/configuration-core/application/portout"
	storageadapterout "go-intelligent-monitoring-system/storage-core/adapterout"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"

	"os"

	"testing"
)

func TestAddAuthorizedFaces(t *testing.T) {

	os.Setenv("AUTHORIZED_FACES_DIRECTORY", "../../test/images/withFaces/")

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = configurationadapterout.NewRekoAdapterTest()

	//Define the "Adapter Out" to be used to connect to the storage core
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = storageadapterout.NewImage2S3AdapterTest()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = configurationapplication.NewFaceIndexerService(storageImageAdapter, rekoAdapter)

	//"Adapter In": DirectoryAdapter gets the authorized faces from a directory
	confDirectoryAdapter := NewDirectoryAdapter(faceIndexerService)

	err := confDirectoryAdapter.AddAuthorizedFaces()

	if err != nil {
		t.Errorf("Error adding authorized faces: %v", err)
	}
}
