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

	for {
		t := time.Now()
		today := t.Format("20060102")

		ftpTodayDirectory := os.Getenv("FTP_DIRECTORY") + today + "/"
		ftpTodayDirectoryProcessed := strings.Replace(ftpTodayDirectory, today, today+"_processed", 1)

		//Create ftpTodayDirectoryProcessed
		_ = os.Mkdir(ftpTodayDirectoryProcessed, os.ModePerm)

		files, err := ioutil.ReadDir(ftpTodayDirectory)
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files {
			var err error
			//TODO: logger.Info("Procesing file: " + f.Name())
			fileBytes, err := ioutil.ReadFile(ftpTodayDirectory + f.Name())
			if err != nil {
				log.Fatal(err)
			}
			if strings.HasSuffix(f.Name(), ".jpg") {
				err = ftp.imageProcessingService.ProcessImage(fileBytes, f.Name())
			} else if strings.HasSuffix(f.Name(), ".mp4") {
				err = ftp.video2ImageService.ProcessVideo(fileBytes, f.Name())
			}

			if err == nil {
				err := os.Rename(ftpTodayDirectory+f.Name(), ftpTodayDirectoryProcessed+f.Name())
				if err != nil {
					log.Fatal(err)
				}
			}

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
