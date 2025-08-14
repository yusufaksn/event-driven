package main

import (
	"log"

	"encoding/json"
	"order/domain"
	"order/infra/couchbase"
	"order/infra/kafka"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	id := uuid.New().String()
	app := fiber.New()
	app.Post("/order", func(c fiber.Ctx) error {
		o := new(domain.OrderItem)
		if err := json.Unmarshal(c.Body(), o); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "invalid JSON"})
		}
		o.OrderID = id
		data, _ := json.Marshal(o)
		orderID := o.OrderID
		kafka.SendKafka(data, orderID)
		couchbase.Save(o, orderID)

		return c.JSON(fiber.Map{"message": "order created"})
	})

	log.Fatal(app.Listen(":3000"))
}
