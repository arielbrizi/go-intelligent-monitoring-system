package storageadapterout

import (
	"context"
	"encoding/json"
	"go-intelligent-monitoring-system/domain"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/segmentio/kafka-go"
)

//KafkaAdapter ...
type KafkaAdapter struct {
	conn   *kafka.Conn
	topic  string
	broker string
}

//SendImage2Queue ...
func (ka *KafkaAdapter) SendImage2Queue(image domain.Image) error {

	value, marshalError := json.Marshal(image)
	if marshalError != nil {
		log.WithFields(log.Fields{"image.Name": image.Name, "image.Bucket": image.Bucket}).WithError(marshalError).Error("Filed to marshal image")
		return marshalError
	}

	i, errWrite := ka.conn.WriteMessages(kafka.Message{Value: value})

	//TODO: retries
	if errWrite != nil {
		log.WithFields(log.Fields{"topic": ka.topic, "broker": ka.broker, "image.Name": image.Name, "image.Bucket": image.Bucket}).WithError(errWrite).Error("Failed to write message to kafka")

		if errClose := ka.conn.Close(); errClose != nil {
			log.WithFields(log.Fields{"image.Name": image.Name, "image.Bucket": image.Bucket}).WithError(errClose).Error("Failed to close writer")
		}

	}

	log.WithFields(log.Fields{"topic": ka.topic, "broker": ka.broker, "bytesWritten": i, "image.Name": image.Name, "image.Bucket": image.Bucket}).Info("Message correctly written to kafka")

	return nil
}

//NewKafkaAdapter initializes an KafkaAdapter object.
func NewKafkaAdapter() *KafkaAdapter {
	// to produce messages
	topic := os.Getenv("QUEUE_TOPIC")
	broker := os.Getenv("QUEUE_BROKER_LIST")
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, partition)
	if err != nil {
		log.WithFields(log.Fields{"topic": topic, "broker": broker}).WithError(err).Fatal("Kafka: failed to dial leader")
	}

	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	return &KafkaAdapter{
		conn:   conn,
		topic:  topic,
		broker: broker,
	}
}
