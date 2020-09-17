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
	imageProcessingService storageapplicationportin.InputImagePort
	video2ImageService     storageapplicationportin.InputVideoPort
}

//Process ...
func (ftp *FtpToInputAdapter) Process() {

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
			ftp.imageProcessingService.ProcessImage(fileBytes)
		} else if strings.HasSuffix(f.Name(), ".mp4") {
			ftp.video2ImageService.ProcessVideo(fileBytes)
		}

	}

}

//NewFtpToInputAdapter initializes an FtpToInputAdapter object.
func NewFtpToInputAdapter(ImageProcessingService storageapplicationportin.InputImagePort, video2ImageService storageapplicationportin.InputVideoPort) *FtpToInputAdapter {
	return &FtpToInputAdapter{
		imageProcessingService: ImageProcessingService,
		video2ImageService:     video2ImageService,
	}
}
