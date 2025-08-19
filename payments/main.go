package main

import (
	"context"

	"log"

	"payments/infra/kafka"
	"payments/infra/mongodb"

	"payments/services"

	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	ctx := context.Background()

	mongoRepo := mongodb.InitMongo(ctx)

	kafkaRepo := kafka.InitKafka()
	defer kafka.CloseKafka(kafkaRepo)

	inventoryService := services.NewPaymementService(kafkaRepo, mongoRepo)

	go inventoryService.ReadKafkaToSaveMongoDB(ctx)

	log.Fatalln(app.Listen(":3002"))
}
