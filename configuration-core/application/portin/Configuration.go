package configurationapplicationportin

import "go-intelligent-monitoring-system/domain"

//ConfigurationPort ...
type ConfigurationPort interface {
	AddAuthorizedFace(image []byte, name string) (*domain.AuthorizedFace, error)
	DeleteAuthorizedFace(domain.AuthorizedFace) error
	GetAuthorizedFaces(collectionName string) ([]domain.AuthorizedFace, error)
}
