package kafka

import (
	"context"
	"fmt"

	"os"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
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

func (r *KafkaRepository) SendKafka(productJson []byte, eventID string, tracer trace.Tracer, ctx context.Context) {

	productItem := kafka.Message{
		Key:   []byte(eventID),
		Value: productJson,
	}

	_, span := tracer.Start(ctx, "kafka publish")

	errWriteMessage := r.writer.WriteMessages(context.Background(), productItem)
	if errWriteMessage != nil {
		fmt.Println("Failed", errWriteMessage)
	} else {
		fmt.Println("The message is sent successfuly")
	}
	span.End()

}

func CloseKafkaConnection(r *KafkaRepository) {
	r.writer.Close()
}
