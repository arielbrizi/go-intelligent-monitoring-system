package recognitionadapterin

import (
	"context"
	"encoding/json"
	"go-intelligent-monitoring-system/domain"
	recognitionapplicationportin "go-intelligent-monitoring-system/recognition-core/application/portin"
	"os"

	log "github.com/sirupsen/logrus"

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
			log.WithError(err).Error("failed to read message from Kafka")
			break
		}

		err = json.Unmarshal(kafkaMessage.Value, &image)
		if err != nil {
			log.WithError(err).Error("failed to Unmarshal image")
		}

		_, err = ka.imageAnalizerService.AnalizeImage(image)
		if err != nil {
			go log.WithFields(log.Fields{"image.Name": image.Name, "kafkaMessage.Topic": kafkaMessage.Topic, "kafkaMessage.Partition": kafkaMessage.Partition, "kafkaMessage.Offset": kafkaMessage.Offset}).Info("Image correctly analized")

		}
	}

	if err = ka.reader.Close(); err != nil {
		log.WithError(err).Fatal("failed to close reader")
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
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaAdapter{
		reader:               r,
		imageAnalizerService: imageAnalizerService,
	}
}
