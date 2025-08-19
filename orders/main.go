package main

import (
	"log"

	"encoding/json"
	"order/domain"
	"order/infra/couchbase"
	"order/infra/kafka"
	"order/services"

	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	kafkaRepo := kafka.InitKafka()
	couchBaseRepo := couchbase.InitCouchBase()
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})
	app.Post("/order", func(c fiber.Ctx) error {

		orderItem := new(domain.OrderItem)
		s := services.NewOrderService(kafkaRepo, couchBaseRepo)

		if err := json.Unmarshal(c.Body(), orderItem); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
		}
		s.SendKafkaToStoreCouchbase(orderItem)

		return c.JSON(fiber.Map{"message": "order created"})
	})
	defer kafka.CloseKafkaConnection(kafkaRepo)
	log.Fatal(app.Listen(":3000"))
}

func IdGenerate() string {
	return uuid.New().String()
}
