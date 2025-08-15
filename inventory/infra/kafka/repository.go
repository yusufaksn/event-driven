package kafka

import (
	"context"
	"encoding/json"
	"inventory/domain"
	"inventory/infra/mongodb"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

var p domain.OrderItem
var reader *kafka.Reader

func InitKafka() {
	brokerAddress := os.Getenv("KAFKA_BROKER_ADDRESS")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP")
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
}

func ReadKafka() {

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error:", err)
			break
		}
		json.Unmarshal(m.Value, &p)
		log.Println("product id :", p.ProductID)
		mongodb.UpdateInventory(p.ProductID, p.Quantity)
	}
}

func CloseKafka() {
	reader.Close()
}
