package kafka

import (
	"context"

	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type KafkaRepository struct {
	reader *kafka.Reader
}

func InitKafka() *KafkaRepository {
	var brokerAddress = os.Getenv("KAFKA_BROKER_ADDRESS")
	var inventory_topic = os.Getenv("KAFKA_INVENTORY_TOPIC")
	var groupID = os.Getenv("KAFKA_GROUP")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		GroupID:        groupID,
		Topic:          inventory_topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	return &KafkaRepository{
		reader: reader,
	}
}

func (r *KafkaRepository) ReadKafka(ctx context.Context) ([]byte, error) {
	m, err := r.reader.ReadMessage(ctx)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	return m.Value, nil
}

func CloseKafka(r *KafkaRepository) {
	r.reader.Close()
}
