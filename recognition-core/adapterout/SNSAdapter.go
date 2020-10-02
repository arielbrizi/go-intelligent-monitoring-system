package recognitionadapterout

import (
	"fmt"
	"go-intelligent-monitoring-system/domain"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

//SNSAdapter ...
type SNSAdapter struct {
	svc *sns.SNS
}

//NotifyTopic ...
func (sn *SNSAdapter) NotifyTopic(notification domain.Notification) error {

	input := &sns.PublishInput{
		Message:  aws.String(notification.Message),
		TopicArn: aws.String(notification.Topic),
	}

	result, err := sn.svc.Publish(input)
	if err != nil {
		fmt.Println("Publish error:", err)
		return err
	}

	fmt.Println(result)

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
