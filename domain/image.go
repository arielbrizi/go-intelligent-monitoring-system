package domain

//Image is the domain struct to represent image attributes to be analized.
type Image struct {
	Name           string `json:"name" binding:"required"`
	Bucket         string `json:"bucket" binding:"required"`
	CollectionName string `json:"collection" binding:"required"`

	URL string `json:"url"` //To get the image file

	Bytes []byte

	//time
	Hour  string
	Day   string
	Month string
	Year  string
}
