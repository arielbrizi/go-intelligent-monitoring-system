package storageapplication

//FilterImageService filter image without faces on it
type FilterImageService struct {
}

//ProcessImage ...
func (filter *FilterImageService) ProcessImage(imgData []byte) (string, error) {
	//TODO
	return "", nil
}

//NewFilterImageService ...
func NewFilterImageService() *FilterImageService {
	return &FilterImageService{}
}
