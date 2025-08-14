package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

/*
*
-- CACHE ASIDE --
*/
var ctx = context.Background()

const (
	host     = "localhost"
	port     = 5432
	user     = "testuser"
	password = "testpass"
	dbname   = "testdb"
)

var brokerAdress = "host.docker.internal:9092"

type productItem struct {
	ProductID   string `json:"product_id"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

func main() {

	productId := "as3d4-ask4-dddc-3337"
	topic := "order_topic"

	msg := productItem{
		ProductID:   productId,
		Name:        "notebook",
		Amount:      15,
		Description: "notebook description",
	}

	productJson, err := json.Marshal(msg)

	if err != nil {
		panic(err)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAdress},
		Topic:   topic,
	})

	defer writer.Close()

	productItem := kafka.Message{
		Key:   []byte(productId),
		Value: productJson,
	}

	errWriteMessage := writer.WriteMessages(context.Background(), productItem)
	if errWriteMessage != nil {
		fmt.Println("Failed", err)
	} else {
		fmt.Println("The message is sent successfuly")
	}

}
