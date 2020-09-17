package domain

//Image is the domain struct to represent image attributes to be analized.
type Image struct {
	Name   string
	Bucket string

	Bytes []byte

	//time
	Hour  string
	Day   string
	Month string
	Year  string
}
