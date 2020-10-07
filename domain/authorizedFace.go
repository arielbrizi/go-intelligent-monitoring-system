package domain

//AuthorizedFace is the domain struct to represent an image wich have an authorized Face to be saved on System
type AuthorizedFace struct {
	Name           string
	Bucket         string
	CollectionName string

	Bytes []byte
}
