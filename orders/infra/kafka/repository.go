package kafka

import (
	"context"
	"fmt"

	"os"

	"github.com/segmentio/kafka-go"
)

type KafkaRepository struct {
	writer *kafka.Writer
}

func InitKafka() *KafkaRepository {
	brokerAddress := os.Getenv("KAFKA_BROKER_ADDRESS")
	topic := os.Getenv("KAFKA_TOPIC")
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})
	return &KafkaRepository{
		writer: writer,
	}
}

func (r *KafkaRepository) SendKafka(productJson []byte, eventID string) {

	productItem := kafka.Message{
		Key:   []byte(eventID),
		Value: productJson,
	}

	errWriteMessage := r.writer.WriteMessages(context.Background(), productItem)
	if errWriteMessage != nil {
		fmt.Println("Failed", errWriteMessage)
	} else {
		fmt.Println("The message is sent successfuly")
	}

}

func CloseKafkaConnection(r *KafkaRepository) {
	r.writer.Close()
}
