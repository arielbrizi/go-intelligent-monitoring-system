package adapterin

import (
	"go-intelligent-monitoring-system/storage-core/application/portin"
)

//FtpToInputAdapter ...
type FtpToInputAdapter struct {
	filterImageService portin.InputPort
}

//ProcessImage ...
func (ftp *FtpToInputAdapter) ProcessImage(image []byte) {
	ftp.filterImageService.ProcessImage(image)
}
