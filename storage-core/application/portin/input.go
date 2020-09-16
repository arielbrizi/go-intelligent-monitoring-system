package storageapplicationportin

//InputImagePort port in for add image to analize
type InputImagePort interface {
	ProcessImage(imgData []byte) (string, error)
}

//InputVideoPort port in for add videos to analize
type InputVideoPort interface {
	ProcessVideo(videoData []byte) ([][]byte, error)
}
