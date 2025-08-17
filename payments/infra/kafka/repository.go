package kafka

import (
	"context"
	"encoding/json"

	"log"
	"os"

	"payments/domain"

	"payments/infra/mongodb"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var reader *kafka.Reader
var inventoryItem domain.InventoryItem
var paymentItem domain.PaymentItem

func InitKafka() {
	var brokerAddress = os.Getenv("KAFKA_BROKER_ADDRESS")
	var inventory_topic = os.Getenv("KAFKA_INVENTORY_TOPIC")
	var groupID = os.Getenv("KAFKA_GROUP")
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		GroupID:        groupID,
		Topic:          inventory_topic,
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
		json.Unmarshal(m.Value, &inventoryItem)
		paymentCalculateToStoreMongo(inventoryItem)

	}
}

func paymentCalculateToStoreMongo(inventoryItem domain.InventoryItem) {
	if inventoryItem.Message == "inventory_reserved" {
		paymentItem.PaymentID = IdGenerate()
		paymentItem.TotalPrice = inventoryItem.Price * float32(inventoryItem.Quantity)
		paymentItem.Message = "payment_received"
		paymentItem.EventID = inventoryItem.EventID
		paymentItem.OrderID = inventoryItem.OrderID
		mongodb.StoreMongoDB(paymentItem)

	} else {
		log.Println("inventroy_failed to rollback..")
	}
}

func IdGenerate() string {
	return uuid.New().String()
}

func CloseKafka() {
	reader.Close()
}
