package storageadapterin

import (
	storageadapterout "go-intelligent-monitoring-system/storage-core/adapterout"
	storageapplication "go-intelligent-monitoring-system/storage-core/application"
	storageapplicationportin "go-intelligent-monitoring-system/storage-core/application/portin"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
	"io/ioutil"
	"time"

	"os"
	"testing"
)

func TestProcess(test *testing.T) {

	os.Setenv("FTP_DIRECTORY", "../../test/images/")
	os.Setenv("SCOPE", "test")

	var t = time.Now()
	var today = t.Format("20060102")

	files, _ := ioutil.ReadDir("../../test/images/withFaces")
	numberOfFiles := len(files)

	//Rename "withFaces" to the pattern used by ftpToInputAdapter, for example: ../../test/images/20201227
	os.Rename("../../test/images/withFaces", "../../test/images/"+today)

	//Define the "Adapter Out" to be used to connect to the Image Queue core: Kafka
	var queueAdapterOut storageapplicationportout.QueueImagePort
	queueAdapterOut = storageadapterout.NewKafkaAdapterTest()

	//Define the "Adapter Out" to be used to connect to the storage core
	var storageImageAdapter storageapplicationportout.StorageImagePort
	storageImageAdapter = storageadapterout.NewImage2S3AdapterTest()

	imageProcessingService := storageapplication.NewImageProcessingService(storageImageAdapter, queueAdapterOut)

	//Define the service to be  used between the "Adapter In" and the "Adapter Out" for Videos
	var video2ImageService storageapplicationportin.InputVideoPort
	video2ImageService = storageapplication.NewVideo2ImageService()

	//"Adapter In": FtpToInputAdapter gets the images to be analized from FTP directory
	ftpToInputAdapter := NewFtpToInputAdapter(imageProcessingService, video2ImageService)
	ftpToInputAdapter.Process()

	//get the number of processed files
	files, _ = ioutil.ReadDir("../../test/images/" + today + "_processed")
	if len(files) != numberOfFiles {
		test.Errorf("Error processing files: shoud be %v and they are %v", numberOfFiles, len(files))
	}

	//rollback directories name
	os.Rename("../../test/images/"+today+"_processed", "../../test/images/withFaces")
	os.Remove("../../test/images/" + today)
	os.Remove("../../test/images/" + today + "_faces_auth")
	os.Remove("../../test/images/" + today + "_faces_not_auth")
	os.Remove("../../test/images/" + today + "_processed")

	//check if rollback is OK
	files, _ = ioutil.ReadDir("../../test/images/withFaces")
	if len(files) != numberOfFiles {
		test.Errorf("Error processing files: shoud be %v and they are %v", numberOfFiles, len(files))
	}

}
