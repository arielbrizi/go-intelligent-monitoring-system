package configurationapplicationportout

import "go-intelligent-monitoring-system/domain"

//ImageRecognitionPort ...
type ImageRecognitionPort interface {
	DeleteCollection(collectionName string) error
	CreateCollection(collectionName string) error
	IndexFace(authorizedFace domain.AuthorizedFace) (*string, error) //Add Authorized Face in Collection
	DeleteFace(authorizedFace domain.AuthorizedFace) error           //Delete indexed Authorized Face in Collection
	ListFaces(collectionName string) ([]domain.AuthorizedFace, error)
}
