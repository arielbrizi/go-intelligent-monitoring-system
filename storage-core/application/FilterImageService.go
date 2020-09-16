package storageapplication

//FilterImageService filter image without faces on it
type FilterImageService struct {
}

//ProcessImage ...
func (filter *FilterImageService) ProcessImage(imgData []byte) (string, error) {
	//TODO Filter images without faces
	//TODO invoke I2Q service
	return "", nil
}

//NewFilterImageService ...
func NewFilterImageService() *FilterImageService {
	return &FilterImageService{}
}
