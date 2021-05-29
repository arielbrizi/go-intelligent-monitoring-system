package storageapplicationportin

//InputImagePort port in for add image to analize
type InputImagePort interface {
	ProcessImage(imgData []byte, fileName string) error
}
