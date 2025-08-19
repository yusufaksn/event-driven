package couchbase

import (
	"log"
	"os"
	"time"

	"github.com/couchbase/gocb/v2"
)

type CouchbaseRepository struct {
	cluster    *gocb.Cluster
	collection *gocb.Collection
}

func InitCouchBase() *CouchbaseRepository {
	var err error
	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_URL"), gocb.ClusterOptions{
		Username: os.Getenv("COUCHBASE_USERNAME"),
		Password: os.Getenv("COUCHBASE_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("Couchbase connection failed: %v", err)
	}
	bucket := cluster.Bucket("default")

	err = bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Fatalf("Bucket not ready: %v", err)
	}
	collection := bucket.DefaultCollection()
	return &CouchbaseRepository{
		cluster:    cluster,
		collection: collection,
	}

}

func (r CouchbaseRepository) Save(data any, orderID string) {
	var errCollection error
	_, errCollection = r.collection.Upsert(orderID, data, nil)
	if errCollection != nil {
		log.Println(errCollection)
	}
	log.Println("Successfully recorded..")
}
