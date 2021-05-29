package recognitionadapterout

import (
	"encoding/json"
	"go-intelligent-monitoring-system/domain"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"

	log "github.com/sirupsen/logrus"
)

//RekoAdapter ...
type RekoAdapter struct {
	bucket string //TODO: just one?
	svc    *rekognition.Rekognition
}

//Recognize ...
func (reko *RekoAdapter) Recognize(image domain.Image) (*domain.AnalizedImage, error) {

	collectionName := image.Bucket //TODO: same as bucket?

	input := &rekognition.SearchFacesByImageInput{
		CollectionId:       aws.String(collectionName),
		FaceMatchThreshold: aws.Float64(80.000000),
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(image.Bucket),
				Name:   aws.String(image.Name),
			},
		},
		MaxFaces: aws.Int64(5),
	}

	result, err := reko.svc.SearchFacesByImage(input)
	if err != nil {
		log.WithFields(log.Fields{"collectionName": collectionName, "image.Name": image.Name, "image.Bucket": image.Bucket}).WithError(err).Error("Error on rekognition.SearchFacesByImageInput")
		return nil, err
	}

	analizedImage, errAnalizedImage := reko.resultToAnalizedImage(result, image)
	if errAnalizedImage != nil {
		log.WithFields(log.Fields{"collectionName": collectionName, "result": result}).WithError(errAnalizedImage).Error("Error generating domain.AnalizedImage")
		return nil, errAnalizedImage
	}

	log.WithFields(log.Fields{"collectionName": collectionName, "result": result}).Info("Image correctly Analized")

	return analizedImage, nil
}

//NewRekoAdapter initializes a RekoAdapter object.
func NewRekoAdapter() *RekoAdapter {

	bucket := os.Getenv("CAMARA_DOMAIN")

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")}, //TODO externalize to environment
	)

	svc := rekognition.New(sess)

	return &RekoAdapter{
		bucket: bucket,
		svc:    svc,
	}
}

func (reko *RekoAdapter) resultToAnalizedImage(result *rekognition.SearchFacesByImageOutput, image domain.Image) (*domain.AnalizedImage, error) {
	var analizedImage domain.AnalizedImage

	analizedImage.Bucket = reko.bucket
	analizedImage.ImageBytes = image.Bytes
	analizedImage.Name = image.Name

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	analizedImage.RecognitionCoreResponse = resultBytes

	//TODO: analize other matches
	if len(result.FaceMatches) > 0 {
		analizedImage.PersonNameDetected = *result.FaceMatches[0].Face.ExternalImageId
	}

	return &analizedImage, nil
}

//////////////////// For Test //////////////////////////

//RekoAdapterTest ...
type RekoAdapterTest struct {
	bucket string //TODO: just one?
	svc    *rekognition.Rekognition
}

//Recognize ...
func (reko *RekoAdapterTest) Recognize(image domain.Image) (*domain.AnalizedImage, error) {
	var analizedImage domain.AnalizedImage
	analizedImage.PersonNameDetected = "Ariel"
	return &analizedImage, nil
}

//NewRekoAdapterTest initializes a RekoAdapterTest object.
func NewRekoAdapterTest() *RekoAdapterTest {

	return &RekoAdapterTest{
		bucket: os.Getenv("CAMARA_DOMAIN"),
	}

}
