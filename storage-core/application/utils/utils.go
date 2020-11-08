package utils

import (
	"bytes"
	"image"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	pigo "github.com/esimov/pigo/core"

	_ "image/jpeg" // It is necessary for pigo implementation.
)

var (
	classifier  *pigo.Pigo
	angle       float64
	cascadeFile []byte
)

func init() {
	var err error
	//---------- Pigo implementation https://github.com/esimov/pigo --------------
	cascadeFile, err = ioutil.ReadFile("config/pigo/facefinder")
	pigo := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err = pigo.Unpack(cascadeFile)
	if err != nil {
		log.Fatalf("Error reading the cascade file: %s", err)
	}
	angle = 0.0 // cascade rotation angle. 0.0 is 0 radians and 1.0 is 2*pi radians
}

//FacesOnImagePigo return the number of faces detected on imgData (using PIGO implementation https://github.com/esimov/pigo)
func FacesOnImagePigo(imgData []byte) (int, error) {
	var err error

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.WithError(err).Error("Error on decode image")
		return 0, err
	}

	src := pigo.ImgToNRGBA(img)

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	dets := classifier.RunCascade(cParams, angle)

	// Calculate the intersection over union (IoU) of two clusters.
	dets = classifier.ClusterDetections(dets, 0.2)

	return len(dets), nil
}

//FacesOnImage return the number of faces detected on imgData
func FacesOnImage(imgData []byte) (int, error) {

	faces, err := FacesOnImagePigo(imgData)

	return faces, err
}
