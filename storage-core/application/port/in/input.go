package in

type inputPort interface {
	processImage(imgData []byte)
	processVideo(videoData []byte)

	processImagePath(imgPath string)
	processVideoPath(videoPath string)
}
