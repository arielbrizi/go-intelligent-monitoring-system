package storageapplication

import (
	storageadapterout "go-intelligent-monitoring-system/storage-core/adapterout"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
	"io/ioutil"

	"testing"
)

func TestProcessImage(t *testing.T) {

	//Define the "Adapter Out" to be used to connect to the Image Queue core: Kafka
	var queueAdapterOut storageapplicationportout.QueueImagePort
	queueAdapterOut = storageadapterout.NewKafkaAdapterTest()

	//Define the "Adapter Out" to be used to connect to the storage core
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = storageadapterout.NewImage2S3AdapterTest()

	imageProcessingService := NewImageProcessingService(storageImageAdapter, queueAdapterOut)

	fileBytes, _ := ioutil.ReadFile("../../test/images/withFaces/3.jpg")
	err := imageProcessingService.ProcessImage(fileBytes, "3.jpg")
	if err != nil {
		t.Errorf("Error processing file: %v", err)
	}

}
