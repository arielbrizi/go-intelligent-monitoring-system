package utils

import (
	"io/ioutil"
	"testing"
)

func TestFacesOnImagePigo(t *testing.T) {
	files, err := ioutil.ReadDir("../../../test/images/withFaces/")
	if err != nil {
		t.Errorf("Error reading files: %v", err)
	}

	for _, f := range files {

		fileBytes, errFile := ioutil.ReadFile("../../../test/images/withFaces/" + f.Name())
		if errFile != nil {
			t.Errorf("Error reading file: %v", errFile)
			return
		}

		i, errorPigo := FacesOnImagePigo(fileBytes)
		if errorPigo != nil {
			t.Errorf("Error analizing image: %v", errorPigo)
		}
		if i < 1 {
			t.Errorf("Face not detected")
		}

	}
}
