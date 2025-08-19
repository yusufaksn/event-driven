package services

import (
	"context"
	"payments/domain"
	"payments/infra/kafka"

	"encoding/json"

	"log"
	"payments/infra/mongodb"

	"github.com/google/uuid"
)

type InventoryService struct {
	kafkaRepo *kafka.KafkaRepository
	mongoRepo *mongodb.MongoRepository
}

func NewInventoryService(kafkaRepo *kafka.KafkaRepository, mongoRepo *mongodb.MongoRepository) *InventoryService {
	return &InventoryService{
		kafkaRepo: kafkaRepo,
		mongoRepo: mongoRepo,
	}
}

func (s *InventoryService) ReadKafkaToSaveMongoDB(ctx context.Context) {
	for {
		var inventoryItem domain.InventoryItem
		resultValue, err := s.kafkaRepo.ReadKafka(ctx)
		if err != nil {
			log.Panicln("Kafka read error:", err)
			continue
		} else {
			json.Unmarshal(resultValue, &inventoryItem)
			s.PaymentCalculateToStoreMongo(inventoryItem)
		}

	}
}

func (s *InventoryService) PaymentCalculateToStoreMongo(inventoryItem domain.InventoryItem) {
	var paymentItem domain.PaymentItem
	if inventoryItem.Message == "inventory_reserved" {
		paymentItem.PaymentID = IdGenerate()
		paymentItem.TotalPrice = inventoryItem.Price * float32(inventoryItem.Quantity)
		paymentItem.Message = "payment_received"
		paymentItem.EventID = inventoryItem.EventID
		paymentItem.OrderID = inventoryItem.OrderID
		s.mongoRepo.StoreMongoDB(paymentItem)

	} else {
		log.Println("inventroy_failed to rollback..")
	}
}

func IdGenerate() string {
	return uuid.New().String()
}
