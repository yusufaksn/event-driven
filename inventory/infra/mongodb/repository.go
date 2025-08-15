package mongodb

import (
	"context"
	"log"
	"time"

	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func InitMongo(ctx context.Context) {
	var err error
	connectionUrl := os.Getenv("MONGO_CONN_URL")
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	if err != nil {
		log.Fatal("MongoDB handles connection errors:", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping hatası:", err)
	}

	log.Println("Connection succesfull..")
}

func UpdateInventory(productID string, totalQuantity int) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		log.Println("Mongo err:", err)
	}

	if res.ModifiedCount == 0 {
		log.Println("Failed.")
	} else {
		log.Println("Success")
	}
}
