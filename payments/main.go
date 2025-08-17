package main

import (
	"context"
	"log"

	"time"

	"payments/infra/kafka"
	"payments/infra/mongodb"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongodb.InitMongo(ctx)
	kafka.InitKafka()
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Concurrency:  256 * 1024,
	})
	kafka.ReadKafka()
	defer kafka.CloseKafka()
	log.Fatalln(app.Listen(":3002"))

}
