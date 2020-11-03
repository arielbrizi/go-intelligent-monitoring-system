package recognitionadapterout

import (
	"go-intelligent-monitoring-system/domain"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//FtpImageStorageAdapter ...
type FtpImageStorageAdapter struct {
	ftpDirectory string
}

//NewFtpImageStorageAdapter initializes a FtpImageStorageAdapter object.
func NewFtpImageStorageAdapter() *FtpImageStorageAdapter {

	return &FtpImageStorageAdapter{
		ftpDirectory: os.Getenv("FTP_DIRECTORY"),
	}

}

//SaveNotAuthorizedImage move image to faces not authorized directory
func (isa *FtpImageStorageAdapter) SaveNotAuthorizedImage(image domain.Image) error {

	var t = time.Now()
	var today = t.Format("20060102")

	var ftpTodayDirectory = isa.ftpDirectory + today + "/"
	var ftpTodayDirectoryProcessed = strings.Replace(ftpTodayDirectory, today, today+"_processed", 1)
	var ftpTodayDirectoryFacesNotAuth = strings.Replace(ftpTodayDirectory, today, today+"_faces_not_auth", 1)

	err := os.Rename(ftpTodayDirectoryProcessed+image.Name, ftpTodayDirectoryFacesNotAuth+image.Name)
	if err != nil {
		log.WithFields(log.Fields{"ftpTodayDirectoryProcessed": ftpTodayDirectoryProcessed, "ftpTodayDirectoryFacesNotAuth": ftpTodayDirectoryFacesNotAuth, "fileName": image.Name}).WithError(err).Error("Error moving file to Not Authorized directory")
	}

	return err
}

//SaveAuthorizedImage move image to faces authorized directory
func (isa *FtpImageStorageAdapter) SaveAuthorizedImage(image domain.Image) error {

	var t = time.Now()
	var today = t.Format("20060102")

	var ftpTodayDirectory = isa.ftpDirectory + today + "/"
	var ftpTodayDirectoryProcessed = strings.Replace(ftpTodayDirectory, today, today+"_processed", 1)
	var ftpTodayDirectoryFacesAuth = strings.Replace(ftpTodayDirectory, today, today+"_faces_auth", 1)

	err := os.Rename(ftpTodayDirectoryProcessed+image.Name, ftpTodayDirectoryFacesAuth+image.Name)
	if err != nil {
		log.WithFields(log.Fields{"ftpTodayDirectoryProcessed": ftpTodayDirectoryProcessed, "ftpTodayDirectoryFacesAuth": ftpTodayDirectoryFacesAuth, "fileName": image.Name}).WithError(err).Error("Error moving file to Authorized directory")
	}

	return err
}
