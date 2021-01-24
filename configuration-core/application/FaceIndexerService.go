package configurationapplication

import (
	configurationapplicationportout "go-intelligent-monitoring-system/configuration-core/application/portout"
	"go-intelligent-monitoring-system/domain"
	"os"

	log "github.com/sirupsen/logrus"
)

//FaceIndexerService manage the images collection
type FaceIndexerService struct {
	recoAdapter configurationapplicationportout.ImageRecognitionPort
}

//AddAuthorizedFace ...
func (fis *FaceIndexerService) AddAuthorizedFace(image []byte, name string) (*domain.AuthorizedFace, error) {

	var authorizedFace domain.AuthorizedFace

	collectionName := os.Getenv("CAMARA_DOMAIN")

	authorizedFace.Name = name
	authorizedFace.Bucket = collectionName
	authorizedFace.Bytes = image
	authorizedFace.CollectionName = collectionName

	faceID, err := fis.recoAdapter.IndexFace(authorizedFace)
	if err != nil {
		return nil, err
	}

	authorizedFace.ID = *faceID

	return &authorizedFace, err
}

//DeleteAuthorizedFace ...
func (fis *FaceIndexerService) DeleteAuthorizedFace(authorizedFace domain.AuthorizedFace) error {

	err := fis.recoAdapter.DeleteFace(authorizedFace)

	return err

}

//GetAuthorizedFaces ...
func (fis *FaceIndexerService) GetAuthorizedFaces(collectionName string) ([]domain.AuthorizedFace, error) {

	authorizedFaces, err := fis.recoAdapter.ListFaces(collectionName)

	return authorizedFaces, err

}

//NewFaceIndexerService ...
func NewFaceIndexerService(recoAdapter configurationapplicationportout.ImageRecognitionPort) *FaceIndexerService {

	fis := &FaceIndexerService{
		recoAdapter: recoAdapter,
	}

	collectionName := os.Getenv("CAMARA_DOMAIN")

	errDel := fis.recoAdapter.DeleteCollection(collectionName)
	if errDel != nil {
		log.WithFields(log.Fields{"collectionName": collectionName}).WithError(errDel).Fatal("Error deleting collection")
	}
	log.WithFields(log.Fields{"collectionName": collectionName}).Info("Collection Deleted")

	errCreate := fis.recoAdapter.CreateCollection(collectionName)
	if errCreate != nil {
		log.WithFields(log.Fields{"collectionName": collectionName}).WithError(errCreate).Fatal("Error Creating collection")
	}
	log.WithFields(log.Fields{"collectionName": collectionName}).Info("Collection Created")

	return fis
}
