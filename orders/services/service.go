package services

import (
	"encoding/json"
	"order/domain"
	"order/infra/couchbase"
	"order/infra/kafka"

	"github.com/google/uuid"
)

type OrderService struct {
	kafkaRepo     *kafka.KafkaRepository
	couchBaseRepo *couchbase.CouchbaseRepository
}

func NewOrderService(kafkaRepo *kafka.KafkaRepository, couchBaseRepo *couchbase.CouchbaseRepository) *OrderService {
	return &OrderService{
		kafkaRepo:     kafkaRepo,
		couchBaseRepo: couchBaseRepo,
	}
}

func (r *OrderService) SendKafkaToStoreCouchbase(orderItem *domain.OrderItem) {
	orderItem.OrderID = idGenerate()
	orderItem.EventID = idGenerate()
	data, _ := json.Marshal(orderItem)
	r.kafkaRepo.SendKafka(data, idGenerate())
	r.couchBaseRepo.Save(orderItem, idGenerate())
}

func idGenerate() string {
	return uuid.New().String()
}
