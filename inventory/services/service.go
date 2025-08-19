package services

import (
	"context"
	"encoding/json"
	"inventory/domain"
	"inventory/infra/kafka"
	"inventory/infra/mongodb"
)

type InventoryService struct {
	kafka *kafka.KafkaRepository
	mongo *mongodb.MongoRepository
}

func NewInventoryService(kafka *kafka.KafkaRepository, mongo *mongodb.MongoRepository) *InventoryService {
	return &InventoryService{
		kafka: kafka,
		mongo: mongo,
	}
}

func (s *InventoryService) ReadKafKaToStoreMongoToPublishKafka(ctx context.Context) {
	var orderItem domain.OrderItem
	for {
		result, _ := s.kafka.ReadKafka(ctx)
		if result != nil {
			json.Unmarshal(result, &orderItem)
			s.mongo.UpdateInventory(orderItem.ProductID, orderItem.Quantity)
			s.publishKafka(orderItem)
		}
	}
}

func (s *InventoryService) publishKafka(orderItem domain.OrderItem) {
	var inventoryItem domain.InventoryItem
	inventoryItem.OrderID = orderItem.OrderID
	inventoryItem.EventID = orderItem.EventID
	inventoryItem.ProductID = orderItem.ProductID
	inventoryItem.Quantity = orderItem.Quantity
	inventoryItem.Message = "inventory_reserved"
	inventoryItem.Price = orderItem.Price
	result, _ := json.Marshal(inventoryItem)
	s.kafka.PublishKafka(result)
}
