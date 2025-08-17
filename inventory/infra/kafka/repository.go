package kafka

import (
	"context"
	"encoding/json"
	"inventory/domain"
	"inventory/infra/mongodb"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var orderItem domain.OrderItem
var inventoryItem domain.InventoryItem
var reader *kafka.Reader
var writer *kafka.Writer

func InitKafka() {
	var brokerAddress = os.Getenv("KAFKA_BROKER_ADDRESS")
	var order_topic = os.Getenv("KAFKA_ORDER_TOPIC")
	var inventory_topic = os.Getenv("KAFKA_INVENTORY_TOPIC")
	var groupID = os.Getenv("KAFKA_GROUP")
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		GroupID:        groupID,
		Topic:          order_topic,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   inventory_topic,
	})
}

func ReadKafka() {
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error:", err)
			break
		}
		json.Unmarshal(m.Value, &orderItem)
		mongodb.UpdateInventory(orderItem.ProductID, orderItem.Quantity)

		publishKafka(orderItem)
	}
}

func publishKafka(orderItem domain.OrderItem) {

	inventoryItem.OrderID = orderItem.OrderID
	inventoryItem.EventID = orderItem.EventID
	inventoryItem.ProductID = orderItem.ProductID
	inventoryItem.Quantity = orderItem.Quantity
	inventoryItem.Message = "inventory_reserved"
	inventoryItem.Price = orderItem.Price
	result, _ := json.Marshal(inventoryItem)

	inventoryData := kafka.Message{
		Key:   []byte(IdGenerate()),
		Value: result,
	}

	errWriteMessage := writer.WriteMessages(context.Background(), inventoryData)
	if errWriteMessage != nil {
		log.Println("Failed", errWriteMessage)
	} else {
		log.Println("The message is sent successfuly")
	}
}

func IdGenerate() string {
	return uuid.New().String()
}

func CloseKafka() {
	reader.Close()
	writer.Close()
}
