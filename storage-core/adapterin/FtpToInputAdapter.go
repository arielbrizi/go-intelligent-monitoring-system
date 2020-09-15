package storageadapterin

import (
	storageapplicationportin "go-intelligent-monitoring-system/storage-core/application/portin"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

//FtpToInputAdapter ...
type FtpToInputAdapter struct {
	filterImageService storageapplicationportin.InputPort
}

//ProcessImage ...
func (ftp *FtpToInputAdapter) processImage(image []byte) {
	ftp.filterImageService.ProcessImage(image)
}

//ProcessImages ...
func (ftp *FtpToInputAdapter) ProcessImages() {

	t := time.Now()
	today := t.Format("20060102")

	ftpTodayDirectory := os.Getenv("FTP_DIRECTORY") + today + "/"

	files, err := ioutil.ReadDir(ftpTodayDirectory)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {

		//TODO: logger.Info("Procesing file: " + f.Name())
		fileBytes, err := ioutil.ReadFile(ftpTodayDirectory + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		if strings.HasSuffix(f.Name(), ".jpg") {
			ftp.filterImageService.ProcessImage(fileBytes)
		} else {
			//ftp.videoToImageService.ProcessVideo(fileBytes)
		}

	}

}

//NewFtpToInputAdapter initializes an FtpToInputAdapter object.
func NewFtpToInputAdapter(filterImageService storageapplicationportin.InputPort) *FtpToInputAdapter {
	return &FtpToInputAdapter{
		filterImageService: filterImageService,
	}
}
