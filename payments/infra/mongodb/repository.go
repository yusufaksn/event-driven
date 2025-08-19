package mongodb

import (
	"context"
	"log"
	"payments/domain"

	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	client *mongo.Client
}

func InitMongo(ctx context.Context) *MongoRepository {
	var err error
	connectionUrl := os.Getenv("MONGO_CONN_URL")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUrl))
	if err != nil {
		log.Fatal("MongoDB handles connection errors:", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping hatası:", err)
	}

	log.Println("Connection succesfull..")
	return &MongoRepository{
		client: client,
	}
}

func (m *MongoRepository) StoreMongoDB(paymentItem domain.PaymentItem) {

	collection := m.client.Database("mydb").Collection("payments")
	filter := bson.M{"orderid": paymentItem.OrderID}
	update := bson.M{"$set": paymentItem}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)

	if err != nil {
		log.Println("Mongo err:", err)
	} else {
		log.Println("Payment record is saved successfully..")
	}

}
