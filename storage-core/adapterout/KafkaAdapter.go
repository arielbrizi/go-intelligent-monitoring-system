package storageadapterout

import (
	"context"
	"encoding/json"
	"go-intelligent-monitoring-system/domain"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/segmentio/kafka-go"
)

//KafkaAdapter ...
type KafkaAdapter struct {
	writer *kafka.Writer
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

	errWrite := ka.writer.WriteMessages(context.Background(), kafka.Message{Value: value})

	if errWrite != nil {

		log.WithFields(log.Fields{"topic": ka.topic, "broker": ka.broker, "image.Name": image.Name, "image.Bucket": image.Bucket}).WithError(errWrite).Error("Failed to write message to kafka")

		ka.writer = kafka.NewWriter(kafka.WriterConfig{
			Brokers:         []string{ka.broker},
			Topic:           ka.topic,
			Balancer:        &kafka.LeastBytes{},
			IdleConnTimeout: 720 * time.Hour, // 1 Mes
		})

		return errWrite
	}

	log.WithFields(log.Fields{"topic": ka.topic, "broker": ka.broker, "image.Name": image.Name, "image.Bucket": image.Bucket}).Info("Message correctly written to kafka")

	return nil
}

//NewKafkaAdapter initializes an KafkaAdapter object.
func NewKafkaAdapter() *KafkaAdapter {
	// to produce messages
	topic := os.Getenv("QUEUE_TOPIC")
	broker := os.Getenv("QUEUE_BROKER_LIST")

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:         strings.Split(broker, ","),
		Topic:           topic,
		Balancer:        &kafka.LeastBytes{},
		IdleConnTimeout: 720 * time.Hour, // 1 Mes
	})

	if w == nil {
		log.WithFields(log.Fields{"topic": topic, "broker": broker}).Fatal("Kafka: failed to create writer")
	} else {
		log.WithFields(log.Fields{"topic": topic, "broker": broker}).Info("Kafka: writer created successfully")

	}

	//createTopic()

	return &KafkaAdapter{
		writer: w,
		topic:  topic,
		broker: broker,
	}
}

///////////////////////// For Test /////////////////

//KafkaAdapterTest ...
type KafkaAdapterTest struct {
	writer *kafka.Writer
	topic  string
	broker string
}

//SendImage2Queue ...
func (ka *KafkaAdapterTest) SendImage2Queue(image domain.Image) error {
	return nil
}

//NewKafkaAdapterTest initializes an KafkaAdapter object.
func NewKafkaAdapterTest() *KafkaAdapterTest {
	ka := NewKafkaAdapter()

	return &KafkaAdapterTest{
		writer: nil,
		topic:  ka.topic,
		broker: ka.broker,
	}
}
