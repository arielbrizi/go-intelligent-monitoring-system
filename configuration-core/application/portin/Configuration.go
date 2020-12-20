package configurationapplicationportin

import "go-intelligent-monitoring-system/domain"

//ConfigurationPort ...
type ConfigurationPort interface {
	AddAuthorizedFace(image []byte, name string) (*domain.AuthorizedFace, error)
	DeleteAuthorizedFace(image []byte, name string) error
}
