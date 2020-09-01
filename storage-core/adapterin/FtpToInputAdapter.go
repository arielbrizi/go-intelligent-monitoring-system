package storageadapterin

import (
	storageapplicationportin "go-intelligent-monitoring-system/storage-core/application/portin"
)

//FtpToInputAdapter ...
type FtpToInputAdapter struct {
	filterImageService storageapplicationportin.InputPort
}

//ProcessImage ...
func (ftp *FtpToInputAdapter) processImage(image []byte) {
	ftp.filterImageService.ProcessImage(image)
}

//ProcessImages ...
func (ftp *FtpToInputAdapter) ProcessImages() {
	//Todo find images from dir
	var image []byte
	/////
	ftp.filterImageService.ProcessImage(image)
}

//NewFtpToInputAdapter initializes an FtpToInputAdapter object.
func NewFtpToInputAdapter(filterImageService storageapplicationportin.InputPort) *FtpToInputAdapter {
	return &FtpToInputAdapter{
		filterImageService: filterImageService,
	}
}
