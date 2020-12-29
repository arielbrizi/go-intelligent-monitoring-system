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

	os.Setenv("FTP_DIRECTORY", "../../test/FtpToInputAdapter/")
	os.Setenv("SCOPE", "test")

	var t = time.Now()
	var today = t.Format("20060102")

	files, _ := ioutil.ReadDir("../../test/FtpToInputAdapter/20201228")
	numberOfFiles := len(files)

	//Change folder name to 'today'
	os.Rename("../../test/FtpToInputAdapter/20201228", "../../test/FtpToInputAdapter/"+today)

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

	files, _ = ioutil.ReadDir("../../test/FtpToInputAdapter/" + today + "_processed")
	if len(files) != numberOfFiles {
		test.Errorf("Error processing files: shoud be %v and they are %v", len(files), numberOfFiles)
	}

	os.Rename("../../test/FtpToInputAdapter/"+today+"_processed", "../../test/FtpToInputAdapter/20201228")
	os.Remove("../../test/FtpToInputAdapter/" + today)
	os.Remove("../../test/FtpToInputAdapter/" + today + "_faces_auth")
	os.Remove("../../test/FtpToInputAdapter/" + today + "_faces_not_auth")
	os.Remove("../../test/FtpToInputAdapter/" + today + "_processed")

	files, _ = ioutil.ReadDir("../../test/FtpToInputAdapter/20201228")
	if len(files) != numberOfFiles {
		test.Errorf("Error processing files: shoud be %v and they are %v", len(files), numberOfFiles)
	}

}
