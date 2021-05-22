package configurationapplication

import (
	configurationapplicationportout "go-intelligent-monitoring-system/configuration-core/application/portout"
	"go-intelligent-monitoring-system/domain"
	storageapplicationportout "go-intelligent-monitoring-system/storage-core/application/portout"
	"os"

	log "github.com/sirupsen/logrus"
)

//FaceIndexerService manage the images collection
type FaceIndexerService struct {
	recoAdapter         configurationapplicationportout.ImageRecognitionPort
	storageImageAdapter storageapplicationportout.StorageImagePort
}

//AddAuthorizedFace ...
func (fis *FaceIndexerService) AddAuthorizedFace(image []byte, name string, bucket string, collectionName string) (*domain.AuthorizedFace, error) {

	var authorizedFace domain.AuthorizedFace

	if collectionName == "" {
		collectionName = os.Getenv("CAMARA_DOMAIN")
	}

	if bucket == "" {
		bucket = collectionName
	}

	if image != nil && len(image) > 1 { //Save image to bucket before indexing it
		err := fis.storageImageAdapter.Save(getImage(image, name, bucket, collectionName))
		if err != nil {
			log.WithFields(log.Fields{"authorizedFace.Name": name, "authorizedFace.Bucket": bucket, "authorizedFace.CollectionName": collectionName}).WithError(err).Error("Error saving authorized face")
			return nil, err
		}
	}

	authorizedFace.Name = name
	authorizedFace.Bucket = bucket
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
func NewFaceIndexerService(storageImageAdapter storageapplicationportout.StorageImagePort, recoAdapter configurationapplicationportout.ImageRecognitionPort) *FaceIndexerService {

	fis := &FaceIndexerService{
		storageImageAdapter: storageImageAdapter,
		recoAdapter:         recoAdapter,
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

func getImage(imageBytes []byte, name string, bucket string, collectionName string) domain.Image {

	var image domain.Image
	image.Bytes = imageBytes
	image.Name = name
	image.Bucket = bucket
	image.CollectionName = collectionName
	return image

}
