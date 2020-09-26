package storageadapterout

import (
	"bytes"
	"fmt"
	"go-intelligent-monitoring-system/domain"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//Image2S3Adapter ...
type Image2S3Adapter struct {
	bucket   string //TODO: just one?
	uploader s3manager.Uploader
}

//Save ...
func (i2s3 *Image2S3Adapter) Save(image domain.Image) error {

	upParams := &s3manager.UploadInput{
		Bucket: &image.Bucket,
		Key:    &image.Name,
		Body:   bytes.NewReader(image.Bytes),
	}

	// Perform an upload.
	_, err := i2s3.uploader.Upload(upParams)

	return err
}

//NewImage2S3Adapter initializes an Image2S3Adapter object.
func NewImage2S3Adapter() *Image2S3Adapter {

	bucket := os.Getenv("CAMARA_DOMAIN")

	// The session the S3 Uploader will use
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// S3 service client the Upload manager will use.
	s3Svc := s3.New(sess)

	// Create an uploader with S3 client and default options
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	// Create the S3 Bucket
	_, err := s3Svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Unable to create bucket %q, %v", bucket, err))
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = s3Svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Error occurred while waiting for bucket to be created, %v", bucket))
	}

	fmt.Printf("Bucket %q successfully created\n", bucket)

	return &Image2S3Adapter{
		bucket:   bucket,
		uploader: *uploader,
	}
}
