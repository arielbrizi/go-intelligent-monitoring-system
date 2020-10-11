package storageadapterin

import (
	storageapplicationportin "go-intelligent-monitoring-system/storage-core/application/portin"
	"io/ioutil"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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
			log.WithFields(log.Fields{"ftpTodayDirectory": ftpTodayDirectory}).WithError(err).Fatal("Error reading directory")
		}

		for _, f := range files {
			var err error

			fileBytes, err := ioutil.ReadFile(ftpTodayDirectory + f.Name())
			if err != nil {
				log.WithFields(log.Fields{"ftpTodayDirectory": ftpTodayDirectory, "fileName": f.Name()}).WithError(err).Fatal("Error reading file")
			}
			if strings.HasSuffix(f.Name(), ".jpg") {
				err = ftp.imageProcessingService.ProcessImage(fileBytes, f.Name())
			} else if strings.HasSuffix(f.Name(), ".mp4") {
				err = ftp.video2ImageService.ProcessVideo(fileBytes, f.Name())
			}

			if err == nil {
				err := os.Rename(ftpTodayDirectory+f.Name(), ftpTodayDirectoryProcessed+f.Name())
				if err != nil {
					log.WithFields(log.Fields{"ftpTodayDirectoryProcessed": ftpTodayDirectoryProcessed, "ftpTodayDirectory": ftpTodayDirectory, "fileName": f.Name()}).WithError(err).Fatal("Error moving file")
				}

				log.WithFields(log.Fields{"ftpTodayDirectoryProcessed": ftpTodayDirectoryProcessed, "ftpTodayDirectory": ftpTodayDirectory, "fileName": f.Name()}).Info("File correctly processed")
			}

		}
	}

}

//NewFtpToInputAdapter initializes an FtpToInputAdapter object.
func NewFtpToInputAdapter(imageProcessingService storageapplicationportin.InputImagePort, video2ImageService storageapplicationportin.InputVideoPort) *FtpToInputAdapter {
	return &FtpToInputAdapter{
		imageProcessingService: imageProcessingService,
		video2ImageService:     video2ImageService,
	}
}
