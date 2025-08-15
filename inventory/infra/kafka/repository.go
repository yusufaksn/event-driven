package kafka

import (
	"context"
	"encoding/json"
	"inventory/infra/mongodb"
	"log"

	"inventory/domain"

	"os"

	"github.com/segmentio/kafka-go"
)

var p domain.OrderItem

func ReadKafka() {
	brokerAddress := os.Getenv("KAFKA_BROKER_ADDRESS")
	topic := os.Getenv("KAFKA_TOPIC")
	groupID := os.Getenv("KAFKA_GROUP")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error:", err)
			break
		}
		json.Unmarshal(m.Value, &p)
		log.Println("product id :", p.ProductID)
		mongodb.UpdateInventory(p.ProductID, p.Quantity)
	}
}
