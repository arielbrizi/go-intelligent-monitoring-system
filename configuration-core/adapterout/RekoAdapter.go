package configurationadapterout

import (
	"go-intelligent-monitoring-system/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	log "github.com/sirupsen/logrus"
)

//RekoAdapter ...
type RekoAdapter struct {
	svc *rekognition.Rekognition
}

//DeleteCollection if exists
func (rekoAdapter *RekoAdapter) DeleteCollection(collectionName string) error {
	input := &rekognition.DeleteCollectionInput{
		CollectionId: aws.String(collectionName),
	}

	result, err := rekoAdapter.svc.DeleteCollection(input)
	if err != nil {
		log.WithFields(log.Fields{"collectionName": collectionName}).WithError(err).Error("Error Deleting collection")
		return err
	}

	log.WithFields(log.Fields{"collectionName": collectionName, "result": result}).Info("Collection deleted")

	return nil
}

//CreateCollection if not exists
func (rekoAdapter *RekoAdapter) CreateCollection(collectionName string) error {

	input := &rekognition.CreateCollectionInput{
		CollectionId: aws.String(collectionName),
	}

	result, err := rekoAdapter.svc.CreateCollection(input)
	if err != nil {
		log.WithFields(log.Fields{"collectionName": collectionName}).WithError(err).Error("Error Creating collection")
		return err
	}

	log.WithFields(log.Fields{"collectionName": collectionName, "result": result}).Info("Collection created")

	return nil
}

//IndexFace add an authorized face to the Collection
func (rekoAdapter *RekoAdapter) IndexFace(image domain.AuthorizedFace) error {

	input := &rekognition.IndexFacesInput{
		CollectionId:    aws.String(image.CollectionName),
		ExternalImageId: aws.String(image.Name),
		Image: &rekognition.Image{
			Bytes: image.Bytes,
		},
	}

	result, err := rekoAdapter.svc.IndexFaces(input)
	if err != nil {
		log.WithFields(log.Fields{"authorizedFace.Name": image.Name, "authorizedFace.Bucket": image.Bucket, "authorizedFace.CollectionName": image.CollectionName, "result": result}).WithError(err).Error("Error indexing face")
		return err
	}

	log.WithFields(log.Fields{"authorizedFace.Name": image.Name, "authorizedFace.Bucket": image.Bucket, "authorizedFace.CollectionName": image.CollectionName, "result": result}).Info("Face indexed")

	return nil
}

//NewRekoAdapter initializes a RekoAdapter object.
func NewRekoAdapter() *RekoAdapter {

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")}, //TODO externalize to environment
	)

	svc := rekognition.New(sess)

	return &RekoAdapter{
		svc: svc,
	}
}
