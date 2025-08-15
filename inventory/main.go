package main

import (
	"context"
	"log"
	"time"

	"inventory/infra/kafka"
	"inventory/infra/mongodb"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongodb.InitMongo(ctx)
	kafka.ReadKafka()
}
