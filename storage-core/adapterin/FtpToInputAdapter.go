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

//Process all images on FTP directory
func (ftp *FtpToInputAdapter) Process() {

	scope := os.Getenv("SCOPE")

	var t = time.Now()
	var today = t.Format("20060102")
	var ftpDirectory = os.Getenv("FTP_DIRECTORY")
	var ftpTodayDirectory = ftpDirectory + today + "/"
	var ftpTodayDirectoryProcessed = strings.Replace(ftpTodayDirectory, today, today+"_processed", 1)
	var ftpTodayDirectoryFacesNotAuth = strings.Replace(ftpTodayDirectory, today, today+"_faces_not_auth", 1)
	var ftpTodayDirectoryFacesAuth = strings.Replace(ftpTodayDirectory, today, today+"_faces_auth", 1)

	errCreateFolders := createFolders(ftpTodayDirectory, ftpTodayDirectoryProcessed, ftpTodayDirectoryFacesNotAuth, ftpTodayDirectoryFacesAuth)
	if errCreateFolders != nil && !strings.Contains(errCreateFolders.Error(), "file exist") {
		log.WithFields(log.Fields{"ftpTodayDirectory": ftpTodayDirectory}).WithError(errCreateFolders).Fatal("Error creating directories")
	}

	for {
		t = time.Now()
		date := t.Format("20060102")

		if date != today { //the day changed
			today = date
			ftpTodayDirectory = ftpDirectory + today + "/"
			ftpTodayDirectoryProcessed = strings.Replace(ftpTodayDirectory, today, today+"_processed", 1)
			ftpTodayDirectoryFacesNotAuth = strings.Replace(ftpTodayDirectory, today, today+"_faces_not_auth", 1)
			ftpTodayDirectoryFacesAuth = strings.Replace(ftpTodayDirectory, today, today+"_faces_auth", 1)

			errCreateFolders := createFolders(ftpTodayDirectory, ftpTodayDirectoryProcessed, ftpTodayDirectoryFacesNotAuth, ftpTodayDirectoryFacesAuth)
			if errCreateFolders != nil && !strings.Contains(errCreateFolders.Error(), "file exist") {
				log.WithFields(log.Fields{"ftpTodayDirectory": ftpTodayDirectory}).WithError(errCreateFolders).Fatal("Error creating directories")
			}
		}

		files, err := ioutil.ReadDir(ftpTodayDirectory)
		if err != nil {
			log.WithFields(log.Fields{"ftpTodayDirectory": ftpTodayDirectory}).WithError(err).Error("Error reading directory")
		}

		for _, f := range files {
			var err error

			fileBytes, errFile := ioutil.ReadFile(ftpTodayDirectory + f.Name())
			if errFile != nil {
				log.WithFields(log.Fields{"ftpTodayDirectory": ftpTodayDirectory, "fileName": f.Name()}).WithError(errFile).Fatal("Error reading file")
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
			} else {
				log.WithFields(log.Fields{"ftpTodayDirectory": ftpTodayDirectory, "fileName": f.Name()}).WithError(err).Error("Error processing file")
			}

		}

		if scope == "test" {
			return
		}
		time.Sleep(5 * time.Second)
	}

}

//NewFtpToInputAdapter initializes an FtpToInputAdapter object.
func NewFtpToInputAdapter(imageProcessingService storageapplicationportin.InputImagePort, video2ImageService storageapplicationportin.InputVideoPort) *FtpToInputAdapter {
	return &FtpToInputAdapter{
		imageProcessingService: imageProcessingService,
		video2ImageService:     video2ImageService,
	}
}

func createFolders(ftpTodayDirectory string, ftpTodayDirectoryProcessed string, ftpTodayDirectoryFacesNotAuth string, ftpTodayDirectoryFacesAuth string) error {
	//Create ftpTodayDirectory
	err := os.Mkdir(ftpTodayDirectory, os.ModePerm)
	if err != nil && !strings.Contains(err.Error(), "file exist") {
		return err
	}

	//Create ftpTodayDirectoryProcessed
	err = os.Mkdir(ftpTodayDirectoryProcessed, os.ModePerm)
	if err != nil && !strings.Contains(err.Error(), "file exist") {
		return err
	}
	//Create ftpTodayDirectoryFacesNotAuth: where images processed are saved.
	err = os.Mkdir(ftpTodayDirectoryFacesNotAuth, os.ModePerm)
	if err != nil && !strings.Contains(err.Error(), "file exist") {
		return err
	}
	//Create ftpTodayDirectoryFacesAuth: where images processed are saved.
	err = os.Mkdir(ftpTodayDirectoryFacesAuth, os.ModePerm)
	if err != nil && !strings.Contains(err.Error(), "file exist") {
		return err
	}

	return err
}
