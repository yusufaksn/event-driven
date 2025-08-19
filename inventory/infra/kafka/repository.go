package kafka

import (
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type KafkaRepository struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func InitKafka() *KafkaRepository {
	var brokerAddress = os.Getenv("KAFKA_BROKER_ADDRESS")
	var order_topic = os.Getenv("KAFKA_ORDER_TOPIC")
	var inventory_topic = os.Getenv("KAFKA_INVENTORY_TOPIC")
	var groupID = os.Getenv("KAFKA_GROUP")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		GroupID:        groupID,
		Topic:          order_topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   inventory_topic,
	})
	return &KafkaRepository{
		reader: reader,
		writer: writer,
	}
}

func (r *KafkaRepository) ReadKafka(ctx context.Context) ([]byte, error) {
	m, err := r.reader.ReadMessage(context.Background())
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	} else {
		return m.Value, nil
	}
}

func (r *KafkaRepository) PublishKafka(data []byte) {
	inventoryData := kafka.Message{
		Key:   []byte(IdGenerate()),
		Value: data,
	}
	errWriteMessage := r.writer.WriteMessages(context.Background(), inventoryData)
	if errWriteMessage != nil {
		log.Println("Failed", errWriteMessage)
	} else {
		log.Println("The message is sent successfuly")
	}
}

func IdGenerate() string {
	return uuid.New().String()
}

func CloseKafka(r *KafkaRepository) {
	r.reader.Close()
	r.writer.Close()
}
