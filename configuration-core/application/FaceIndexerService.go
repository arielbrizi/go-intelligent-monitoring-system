package configurationapplication

import (
	configurationapplicationportout "go-intelligent-monitoring-system/configuration-core/application/portout"
	"go-intelligent-monitoring-system/domain"
	"log"
	"os"
)

//FaceIndexerService manage the images collection
type FaceIndexerService struct {
	rekoAdapter configurationapplicationportout.ImageRecognitionPort
}

//AddAuthorizedFace ...
func (fis *FaceIndexerService) AddAuthorizedFace(image []byte, name string) error {

	var authorizedFace domain.AuthorizedFace
	authorizedFace.Name = name
	authorizedFace.Bucket = os.Getenv("CAMARA_DOMAIN")
	authorizedFace.Bytes = image
	authorizedFace.CollectionName = os.Getenv("CAMARA_DOMAIN")

	fis.rekoAdapter.IndexFace(authorizedFace)

	return nil
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

	fis.rekoAdapter.DeleteCollection(os.Getenv("CAMARA_DOMAIN"))

	err := fis.rekoAdapter.CreateCollection(os.Getenv("CAMARA_DOMAIN"))

	if err != nil {
		log.Fatal(err)
	}

	return fis
}
