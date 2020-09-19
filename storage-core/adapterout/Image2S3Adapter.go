package storageadapterout

import (
	"go-intelligent-monitoring-system/domain"
)

//Image2S3Adapter ...
type Image2S3Adapter struct {
}

//Save ...
func (i2s3 *Image2S3Adapter) Save(image domain.Image) error {

	//TODO
	return nil
}

//NewImage2S3Adapter initializes an Image2S3Adapter object.
func NewImage2S3Adapter() *Image2S3Adapter {
	return &Image2S3Adapter{

	}
}
