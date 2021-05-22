package configurationadapterout

import (
	"errors"
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
func (rekoAdapter *RekoAdapter) IndexFace(image domain.AuthorizedFace) (*string, error) {

	input := &rekognition.IndexFacesInput{
		CollectionId:    aws.String(image.CollectionName),
		ExternalImageId: aws.String(image.Name),
		MaxFaces:        aws.Int64(1), // Just index the biggest face
		Image: &rekognition.Image{
			S3Object: &rekognition.S3Object{
				Bucket: aws.String(image.Bucket),
				Name:   aws.String(image.Name),
			},
		},
	}

	result, err := rekoAdapter.svc.IndexFaces(input)

	if err != nil {
		log.WithFields(log.Fields{"authorizedFace.Name": image.Name, "authorizedFace.Bucket": image.Bucket, "authorizedFace.CollectionName": image.CollectionName, "result": result}).WithError(err).Error("Error indexing face")
		return nil, err
	}

	log.WithFields(log.Fields{"authorizedFace.Name": image.Name, "authorizedFace.Bucket": image.Bucket, "authorizedFace.CollectionName": image.CollectionName, "result": result}).Info("Face indexed")

	if len(result.FaceRecords) < 1 {
		log.WithFields(log.Fields{"authorizedFace.Name": image.Name, "authorizedFace.Bucket": image.Bucket, "authorizedFace.CollectionName": image.CollectionName, "result": result}).Error("No face detected")
		return nil, errors.New("No face detected")
	}

	return result.FaceRecords[0].Face.FaceId, nil //[0] becouse MaxFaces was set to 1.
}

//DeleteFace delete an indexed authorized face in the Collection
func (rekoAdapter *RekoAdapter) DeleteFace(authorizedFace domain.AuthorizedFace) error {

	input := &rekognition.DeleteFacesInput{
		CollectionId: aws.String(authorizedFace.CollectionName),
		FaceIds: []*string{
			aws.String(authorizedFace.ID),
		},
	}

	result, err := rekoAdapter.svc.DeleteFaces(input)
	if err != nil {
		log.WithFields(log.Fields{"authorizedFace.ID": authorizedFace.ID, "result": result}).WithError(err).Error("Error deleting face")
		return err
	}

	log.WithFields(log.Fields{"authorizedFace.Name": authorizedFace.ID, "result": result}).Info("Face deleted")

	return nil
}

//ListFaces get indexed authorized faces in a Collection
func (rekoAdapter *RekoAdapter) ListFaces(collectionName string) ([]domain.AuthorizedFace, error) {

	var authorizedFaces []domain.AuthorizedFace

	input := &rekognition.ListFacesInput{
		CollectionId: aws.String(collectionName),
		//MaxResults:   aws.Int64(20),
	}

	result, err := rekoAdapter.svc.ListFaces(input)
	if err != nil {
		log.WithFields(log.Fields{"collectionName": collectionName, "result": result}).WithError(err).Error("Error listing faces")
		return nil, err
	}

	log.WithFields(log.Fields{"collectionName": collectionName, "result": result}).Info("Indexed faces: ", len(result.Faces))

	for _, face := range result.Faces {
		var authorizedFace domain.AuthorizedFace
		authorizedFace.Name = *face.ExternalImageId //It was set by me when It was indexed
		authorizedFace.CollectionName = collectionName
		authorizedFace.ID = *face.FaceId
		authorizedFaces = append(authorizedFaces, authorizedFace)
	}

	return authorizedFaces, nil
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

////////////////// For Test ////////////////////////

//RekoAdapterTest ...
type RekoAdapterTest struct {
}

//NewRekoAdapterTest initializes a RekoAdapter object.
func NewRekoAdapterTest() *RekoAdapterTest {
	return &RekoAdapterTest{}
}

//IndexFace add an authorized face to the Collection
func (rekoAdapter *RekoAdapterTest) IndexFace(image domain.AuthorizedFace) (*string, error) {
	var faceID = "123"
	return &faceID, nil
}

//DeleteFace delete an indexed authorized face in the Collection
func (rekoAdapter *RekoAdapterTest) DeleteFace(image domain.AuthorizedFace) error {
	return nil
}

//ListFaces get indexed authorized faces in the Collection
func (rekoAdapter *RekoAdapterTest) ListFaces(collectionName string) ([]domain.AuthorizedFace, error) {
	var authorizedFaces []domain.AuthorizedFace
	var authorizedFace domain.AuthorizedFace
	authorizedFace.Name = "silvia1.jpg"
	authorizedFace.Bucket = "camarasilvia"
	authorizedFace.CollectionName = "camarasilvia"
	authorizedFace.ID = "2659022e-4ad2-4be8-81b9-1d4b1953ff90"
	authorizedFaces = append(authorizedFaces, authorizedFace)
	return authorizedFaces, nil
}

//CreateCollection if not exists
func (rekoAdapter *RekoAdapterTest) CreateCollection(collectionName string) error {
	return nil
}

//DeleteCollection if exists
func (rekoAdapter *RekoAdapterTest) DeleteCollection(collectionName string) error {
	return nil
}
