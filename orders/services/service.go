package services

import (
	"encoding/json"
	"order/domain"
	"order/infra/couchbase"
	"order/infra/kafka"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type OrderService struct {
	kafkaRepo     *kafka.KafkaRepository
	couchBaseRepo *couchbase.CouchbaseRepository
	Tracer        trace.Tracer
}

func NewOrderService(kafkaRepo *kafka.KafkaRepository, couchBaseRepo *couchbase.CouchbaseRepository) *OrderService {
	return &OrderService{
		kafkaRepo:     kafkaRepo,
		couchBaseRepo: couchBaseRepo,
	}
}

func (r *OrderService) SendKafkaToStoreCouchbase(orderItem *domain.OrderItem, tracer trace.Tracer) {
	orderItem.OrderID = idGenerate()
	orderItem.EventID = idGenerate()
	data, _ := json.Marshal(orderItem)
	ctx := r.couchBaseRepo.Save(orderItem, idGenerate(), tracer)
	r.kafkaRepo.SendKafka(data, idGenerate(), tracer, ctx)

}

func idGenerate() string {
	return uuid.New().String()
}
