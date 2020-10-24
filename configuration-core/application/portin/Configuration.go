package configurationapplicationportin

//ConfigurationPort ...
type ConfigurationPort interface {
	AddAuthorizedFace(image []byte, name string) error
	DeleteAuthorizedFace(image []byte, name string) error
}
