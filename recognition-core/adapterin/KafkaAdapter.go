package recognitionadapterin

import (
	"context"
	"encoding/json"
	"go-intelligent-monitoring-system/domain"
	recognitionapplicationportin "go-intelligent-monitoring-system/recognition-core/application/portin"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

//KafkaAdapter ...
type KafkaAdapter struct {
	reader               *kafka.Reader
	imageAnalizerService recognitionapplicationportin.QueueImagePort
}

//ReceiveImagesFromQueue ...
func (ka *KafkaAdapter) ReceiveImagesFromQueue() error {
	var err error
	var kafkaMessage kafka.Message
	var image domain.Image
	for {
		kafkaMessage, err = ka.reader.ReadMessage(context.Background())

		if err != nil {
			break
		}

		err = json.Unmarshal(kafkaMessage.Value, &image)
		if err != nil {
			log.Fatal("failed to Unmarshal image:", err)
		}

		err = ka.imageAnalizerService.AnalizeImage(image)
		if err != nil {
			log.Fatal("failed to analize image:", err)
		}

		//fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", kafkaMessage.Topic, kafkaMessage.Partition, kafkaMessage.Offset, string(kafkaMessage.Key), string(kafkaMessage.Value))
	}

	if err = ka.reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	return err

}

//NewKafkaAdapter initializes an KafkaAdapter object.
func NewKafkaAdapter(imageAnalizerService recognitionapplicationportin.QueueImagePort) *KafkaAdapter {
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
		reader:               r,
		imageAnalizerService: imageAnalizerService,
	}
}
