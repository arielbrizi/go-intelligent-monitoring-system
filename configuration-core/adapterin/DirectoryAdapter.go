package configurationadapterin

import (
	configurationapplicationportin "go-intelligent-monitoring-system/configuration-core/application/portin"
	"io/ioutil"
	"log"
	"os"
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
		log.Fatal(err)
	}

	for _, f := range files {
		var err error
		//TODO: logger.Info("Procesing file: " + f.Name())
		fileBytes, err := ioutil.ReadFile(directory + f.Name())
		if err != nil {
			log.Fatal(err)
		}

		err = da.faceIndexerService.AddAuthorizedFace(fileBytes, f.Name())
		if err != nil {
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
