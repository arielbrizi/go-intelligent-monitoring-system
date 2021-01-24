package domain

//AuthorizedFace is the domain struct to represent an image wich have an authorized Face to be saved on System
type AuthorizedFace struct {
	Name           string `json:"name"`
	Bucket         string `json:"bucket"`
	CollectionName string `json:"collection" binding:"required"`
	ID             string `json:"id" binding:"required"`

	Bytes []byte
}
