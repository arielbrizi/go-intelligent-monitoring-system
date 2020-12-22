package configurationadapterin

import (
	configurationadapterout "go-intelligent-monitoring-system/configuration-core/adapterout"
	configurationapplication "go-intelligent-monitoring-system/configuration-core/application"
	configurationapplicationportin "go-intelligent-monitoring-system/configuration-core/application/portin"
	configurationapplicationportout "go-intelligent-monitoring-system/configuration-core/application/portout"
	"os"

	"testing"
)

func TestAddAuthorizedFaces(t *testing.T) {

	os.Setenv("AUTHORIZED_FACES_DIRECTORY", "../../test/images/withFaces/")

	//Define the "Adapter Out" to be used to connect to the recognition core
	var rekoAdapter configurationapplicationportout.ImageRecognitionPort
	rekoAdapter = configurationadapterout.NewRekoAdapterTest()

	//Define the service to be  used between the "Adapter In" and the "Adapter Out"
	var faceIndexerService configurationapplicationportin.ConfigurationPort
	faceIndexerService = configurationapplication.NewFaceIndexerService(rekoAdapter)

	//"Adapter In": DirectoryAdapter gets the authorized faces from a directory
	confDirectoryAdapter := NewDirectoryAdapter(faceIndexerService)

	err := confDirectoryAdapter.AddAuthorizedFaces()

	if err != nil {
		t.Errorf("Error adding authorized faces: %v", err)
	}
}
