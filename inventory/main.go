package main

import (
	"context"
	"log"
	"time"

	"inventory/infra/kafka"
	"inventory/infra/mongodb"
	"inventory/services"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	kafkaRepo := kafka.InitKafka()
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoRepo := mongodb.InitMongo(ctx)
	s := services.NewInventoryService(kafkaRepo, mongoRepo)
	s.ReadKafKaToStoreMongoToPublishKafka(ctx)

	defer kafka.CloseKafka(kafkaRepo)

	log.Fatalln(app.Listen(":3001"))
}
