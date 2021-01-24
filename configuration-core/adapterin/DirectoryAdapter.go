package configurationadapterin

import (
	configurationapplicationportin "go-intelligent-monitoring-system/configuration-core/application/portin"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

//DirectoryAdapter represents the adapter wich adds to the monitoring system  all authorized faces located on defined images directory
type DirectoryAdapter struct {
	faceIndexerService configurationapplicationportin.ConfigurationPort
}

//AddAuthorizedFaces ...
func (da *DirectoryAdapter) AddAuthorizedFaces() error {

	directory := os.Getenv("AUTHORIZED_FACES_DIRECTORY")

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.WithFields(log.Fields{"directory": directory}).WithError(err).Fatal("Reading AUTHORIZED_FACES_DIRECTORY directory")
	}

	for _, f := range files {
		var err error

		log.WithFields(log.Fields{"file": f.Name()}).Info("Procesing file")

		fileBytes, err := ioutil.ReadFile(directory + f.Name())
		if err != nil {
			log.WithFields(log.Fields{"file": f.Name()}).WithError(err).Fatal("Procesing file")
		}

		//TODO: save authorized Faces to be able to delete some of them (with faceID)
		_, err = da.faceIndexerService.AddAuthorizedFace(fileBytes, f.Name())
		if err != nil {
			log.WithFields(log.Fields{"file": f.Name()}).WithError(err).Error("Adding authorized face")
			return err
		}
	}

	return nil

}

//NewDirectoryAdapter initializes an DirectoryAdapter object.
func NewDirectoryAdapter(faceIndexerService configurationapplicationportin.ConfigurationPort) *DirectoryAdapter {
	return &DirectoryAdapter{
		faceIndexerService: faceIndexerService,
	}
}
