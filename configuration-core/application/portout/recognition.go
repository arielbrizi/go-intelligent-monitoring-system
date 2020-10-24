package configurationapplicationportout

import "go-intelligent-monitoring-system/domain"

//ImageRecognitionPort ...
type ImageRecognitionPort interface {
	DeleteCollection(collectionName string) error
	CreateCollection(collectionName string) error
	IndexFace(authorizedFace domain.AuthorizedFace) error //Add Authorized Face in Collection
}
