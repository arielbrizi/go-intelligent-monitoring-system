package recognitionadapterin

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

//KafkaAdapter ...
type KafkaAdapter struct {
	reader *kafka.Reader
}

//ReceiveImagesFromQueue ...
func (ka *KafkaAdapter) ReceiveImagesFromQueue() error {
	var err error
	var kafkaMessage kafka.Message
	for {
		kafkaMessage, err = ka.reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		// TODO: unmarsahl to Image
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", kafkaMessage.Topic, kafkaMessage.Partition, kafkaMessage.Offset, string(kafkaMessage.Key), string(kafkaMessage.Value))
	}

	if err = ka.reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	return err

}

//NewKafkaAdapter initializes an KafkaAdapter object.
func NewKafkaAdapter() *KafkaAdapter {
	// to produce messages
	topic := os.Getenv("QUEUE_TOPIC")
	broker := os.Getenv("QUEUE_BROKER_LIST")

	// make a new reader that consumes from topic-A
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  "consumer-group-id",
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaAdapter{
		reader: r,
	}
}
