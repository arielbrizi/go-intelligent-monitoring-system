package storageadapterout

import (
	"go-intelligent-monitoring-system/domain"
)

//KafkaAdapter ...
type KafkaAdapter struct {
}

//SendImage2Queue ...
func (ka *KafkaAdapter) SendImage2Queue(image domain.Image) error {

	//TODO SendImage2Queue

	return nil
}

//NewKafkaAdapter initializes an KafkaAdapter object.
func NewKafkaAdapter() *KafkaAdapter {
	return &KafkaAdapter{

	}
}
