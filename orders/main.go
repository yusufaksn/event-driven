package main

import (
	"log"

	"encoding/json"
	"order/domain"
	"order/infra/couchbase"
	"order/infra/kafka"

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
	kafka.InitKafka()
	couchbase.InitCouchBase()
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})
	app.Post("/order", func(c fiber.Ctx) error {
		o := new(domain.OrderItem)
		if err := json.Unmarshal(c.Body(), o); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
		}
		o.OrderID = IdGenerate()
		o.EventID = IdGenerate()
		data, _ := json.Marshal(o)
		kafka.SendKafka(data, IdGenerate())
		couchbase.Save(o, IdGenerate())

		return c.JSON(fiber.Map{"message": "order created"})
	})
	defer kafka.CloseKafkaConnection()
	log.Fatal(app.Listen(":3000"))
}

func IdGenerate() string {
	return uuid.New().String()
}
