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

//NotifyTopic ...
func (sn *SNSAdapter) NotifyTopic(notification domain.Notification) error {

	url := "https://" + notification.Image.Bucket + ".s3.amazonaws.com/" + notification.Image.Name

	msg := notification.Message + "\n \n " + url + "\n \n "

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

//NotifySMS ...
func (sn *SNSAdapter) NotifySMS(notification domain.SMSNotification) error {
	//NOT USED ON AWS. THE TYPE OF CHANNEL IT'S CONFIGURED ON TOPIC
	return nil
}

//NotifyEmail ...
func (sn *SNSAdapter) NotifyEmail(notification domain.EmailNotification) error {
	//NOT USED ON AWS. THE TYPE OF CHANNEL IT'S CONFIGURED ON TOPIC
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
