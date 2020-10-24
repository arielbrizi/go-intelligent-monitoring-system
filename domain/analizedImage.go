package domain

//AnalizedImage is the domain struct to represent image attributes analized.
type AnalizedImage struct {
	Name   string
	Bucket string

	PersonNameDetected string

	ImageBytes []byte

	RecognitionCoreResponse []byte //json

	//time
	Hour  string
	Day   string
	Month string
	Year  string
}
