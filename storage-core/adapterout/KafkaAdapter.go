package storageadapterout

import (
	"context"
	"encoding/json"
	"fmt"
	"go-intelligent-monitoring-system/domain"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

//KafkaAdapter ...
type KafkaAdapter struct {
	conn *kafka.Conn
}

//SendImage2Queue ...
func (ka *KafkaAdapter) SendImage2Queue(image domain.Image) error {

	value, marshalError := json.Marshal(image)
	if marshalError != nil {
		log.Fatal(fmt.Errorf("Failed to marshal image"))
	}

	i, errWrite := ka.conn.WriteMessages(kafka.Message{Value: value})

	//TODO: retries
	if errWrite != nil {
		if errClose := ka.conn.Close(); errClose != nil {
			fmt.Printf("failed to close writer: %s \n", errClose)
		}
		log.Fatal(fmt.Errorf("failed to write messages: %v", errWrite))
	}

	fmt.Printf("%v bytes written \n", i)
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
		log.Fatal("failed to dial leader:", err)
	}

	//conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

	return &KafkaAdapter{
		conn: conn,
	}
}
