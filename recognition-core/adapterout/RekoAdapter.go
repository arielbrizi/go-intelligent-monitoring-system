package recognitionadapterout

import (
	"encoding/json"
	"fmt"
	"go-intelligent-monitoring-system/domain"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

//RekoAdapter ...
type RekoAdapter struct {
	bucket string //TODO: just one?
	svc    *rekognition.Rekognition
}

//Recognize ...
func (reko *RekoAdapter) Recognize(image domain.Image) (*domain.AnalizedImage, error) {

	input := &rekognition.SearchFacesByImageInput{
		CollectionId:       aws.String(image.Bucket), //TODO: same as bucket?
		FaceMatchThreshold: aws.Float64(95.000000),
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
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case rekognition.ErrCodeInvalidS3ObjectException:
				fmt.Println(rekognition.ErrCodeInvalidS3ObjectException, aerr.Error())
			case rekognition.ErrCodeInvalidParameterException:
				fmt.Println(rekognition.ErrCodeInvalidParameterException, aerr.Error())
			case rekognition.ErrCodeImageTooLargeException:
				fmt.Println(rekognition.ErrCodeImageTooLargeException, aerr.Error())
			case rekognition.ErrCodeAccessDeniedException:
				fmt.Println(rekognition.ErrCodeAccessDeniedException, aerr.Error())
			case rekognition.ErrCodeInternalServerError:
				fmt.Println(rekognition.ErrCodeInternalServerError, aerr.Error())
			case rekognition.ErrCodeThrottlingException:
				fmt.Println(rekognition.ErrCodeThrottlingException, aerr.Error())
			case rekognition.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(rekognition.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case rekognition.ErrCodeResourceNotFoundException:
				fmt.Println(rekognition.ErrCodeResourceNotFoundException, aerr.Error())
			case rekognition.ErrCodeInvalidImageFormatException:
				fmt.Println(rekognition.ErrCodeInvalidImageFormatException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, err
	}

	analizedImage, errAnalizedImage := reko.resultToAnalizedImage(result, image)
	if errAnalizedImage != nil {
		return nil, errAnalizedImage
	}

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
