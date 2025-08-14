package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"encoding/json"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

var uri = "mongodb://admin:secret123@localhost:27017/mydb?authSource=admin"

func initMongo(ctx context.Context) {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("MongoDB bağlantı hatası:", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping hatası:", err)
	}

	fmt.Println("MongoDB bağlantısı başarılı!")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ReadKafka()

	/*productId := "484bed38-7709-415a-a447-bc3d77138d8a"
	UpdateInventory(productId, 2)*/
}

type OrderItem struct {
	ProductID   string `json:"productID"`
	OrderID     string `json:"orderID"`
	Quantity    int    `json:"quantity"`
	Total       int    `json:"total"`
	Description string `json:"description"`
}

func ReadKafka() {
	brokerAddress := "host.docker.internal:9092"
	topic := "order_topic"
	groupID := "my-group"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{brokerAddress},
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: 0,    // Manuel commit için
	})
	defer r.Close()
	var p OrderItem
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Mesaj okuma hatası:", err)
			break
		}
		/*json.Unmarshal(m.Value, &p)

		log.Printf("result :",  p.ProductID)*/

		json.Unmarshal(m.Value, &p)
		log.Println("product id :", p.ProductID)
		UpdateInventory(p.ProductID, p.Quantity)

	}
}

func UpdateInventory(productID string, totalQuantity int) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	initMongo(ctx)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("MongoDB disconnection..", err)
		}
	}()

	collection := client.Database("mydb").Collection("mydb")

	filter := bson.D{
		{Key: "productID", Value: productID},
		{Key: "total_quantity", Value: bson.D{{Key: "$gt", Value: 0}}},
	}

	update := bson.D{
		{Key: "$inc", Value: bson.D{{Key: "total_quantity", Value: -totalQuantity}}},
	}

	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if res.ModifiedCount == 0 {
		fmt.Println("Failed.")
	} else {
		fmt.Println("Success")
	}
}
