package configurationapplication

import (
	configurationapplicationportout "go-intelligent-monitoring-system/configuration-core/application/portout"
	"go-intelligent-monitoring-system/domain"
	"os"

	log "github.com/sirupsen/logrus"
)

//FaceIndexerService manage the images collection
type FaceIndexerService struct {
	rekoAdapter configurationapplicationportout.ImageRecognitionPort
}

//AddAuthorizedFace ...
func (fis *FaceIndexerService) AddAuthorizedFace(image []byte, name string) (*domain.AuthorizedFace, error) {

	var authorizedFace domain.AuthorizedFace

	collectionName := os.Getenv("CAMARA_DOMAIN")

	authorizedFace.Name = name
	authorizedFace.Bucket = collectionName
	authorizedFace.Bytes = image
	authorizedFace.CollectionName = collectionName

	err := fis.rekoAdapter.IndexFace(authorizedFace)
	if err != nil {
		return nil, err
	}

	return &authorizedFace, err
}

//DeleteAuthorizedFace ...
func (fis *FaceIndexerService) DeleteAuthorizedFace(image []byte, name string) error {

	//TODO DeleteAuthorizedFace

	return nil
}

//NewFaceIndexerService ...
func NewFaceIndexerService(rekoAdapter configurationapplicationportout.ImageRecognitionPort) *FaceIndexerService {

	fis := &FaceIndexerService{
		rekoAdapter: rekoAdapter,
	}

	collectionName := os.Getenv("CAMARA_DOMAIN")

	fis.rekoAdapter.DeleteCollection(collectionName)

	err := fis.rekoAdapter.CreateCollection(collectionName)

	if err != nil {
		log.WithFields(log.Fields{"collectionName": collectionName}).WithError(err).Fatal("Error Creating collection")
	}

	return fis
}
