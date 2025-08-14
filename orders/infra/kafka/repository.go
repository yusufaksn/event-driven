package kafka

import (
	"context"
	"fmt"

	"os"

	"github.com/segmentio/kafka-go"
)

func SendKafka(productJson []byte, productId string) {
	brokerAddress := os.Getenv("KAFKA_BROKER_ADDRESS")
	topic := os.Getenv("KAFKA_TOPIC")
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	defer writer.Close()

	productItem := kafka.Message{
		Key:   []byte(productId),
		Value: productJson,
	}

	errWriteMessage := writer.WriteMessages(context.Background(), productItem)
	if errWriteMessage != nil {
		fmt.Println("Failed", errWriteMessage)
	} else {
		fmt.Println("The message is sent successfuly")
	}

}
