package main

import (
	"context"
	"log"
	"time"

	"inventory/infra/kafka"
	"inventory/infra/mongodb"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	kafka.InitKafka()
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongodb.InitMongo(ctx)
	kafka.ReadKafka()
	defer kafka.CloseKafka()
	log.Fatal(app.Listen(":3001"))
}
