package domain

//Image is the domain struct to represent image attributes.
type Image struct {
	name   string
	bucket string

	//time
	hour  string
	day   string
	month string
	year  string
}
