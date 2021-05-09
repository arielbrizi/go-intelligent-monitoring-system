package recognitionadapterout

import (
	"go-intelligent-monitoring-system/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	log "github.com/sirupsen/logrus"
)

//SNSAdapter ...
type SNSAdapter struct {
	svc *sns.SNS
}

//NotifyInitializedSystem
func (ta *SNSAdapter) NotifyInitializedSystem() error {
	//TODO
	return nil
}

//NotifyUnauthorizedFace ...
func (sn *SNSAdapter) NotifyUnauthorizedFace(notification domain.Notification) error {

	msg := notification.Message + "\n \n " + notification.Image.URL + "\n \n "

	input := &sns.PublishInput{
		Message:  aws.String(msg),
		TopicArn: aws.String(notification.Topic),
	}

	result, err := sn.svc.Publish(input)
	if err != nil {
		log.WithFields(log.Fields{"notification.Topic": notification.Topic, "notification.Message": notification.Message, "result": result}).WithError(err).Error("Error on publishing message")
		return err
	}

	log.WithFields(log.Fields{"notification.Topic": notification.Topic, "notification.Message": notification.Message, "result": result}).Info("Message correctly sent")

	return nil
}

//NewSNSAdapter initializes a SNSAdapter object.
func NewSNSAdapter() *SNSAdapter {

	mySession := session.Must(session.NewSession())

	//TODO: create Topic and Suscriptors

	// Create a SNS client with additional configuration
	svc := sns.New(mySession, aws.NewConfig().WithRegion("us-east-1"))

	return &SNSAdapter{
		svc: svc,
	}

}

///////////////////// For Test //////////////////////////////////////

//NewSNSAdapterTest initializes a SNSAdapterTest object.
func NewSNSAdapterTest() *SNSAdapterTest {

	return &SNSAdapterTest{}

}

//SNSAdapterTest ...
type SNSAdapterTest struct {
}

//NotifyTopic ...
func (sn *SNSAdapterTest) NotifyUnauthorizedFace(notification domain.Notification) error {
	return nil
}

//NotifyInitializedSystem ...
func (sn *SNSAdapterTest) NotifyInitializedSystem() error {
	return nil
}
