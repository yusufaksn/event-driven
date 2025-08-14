package couchbase

import (
	"log"
	"os"

	"github.com/couchbase/gocb/v2"
)

func Save(data any, orderID string) {

	cluster, err := gocb.Connect(os.Getenv("COUCHBASE_URL"), gocb.ClusterOptions{
		Username: os.Getenv("COUCHBASE_USERNAME"),
		Password: os.Getenv("COUCHBASE_PASSWORD"),
	})
	if err != nil {
		log.Fatalln(err)

	}
	/*
		bucket := cluster.Bucket("default")
		operation := func() error {
			err := bucket.WaitUntilReady(5*time.Second, nil)
			if err != nil {
				log.Println("Bucket not ready yet")
			}
			return err
		}
		// Exponential backoff ayarları
		b := backoff.NewExponentialBackOff()
		b.MaxElapsedTime = 15 * time.Second // Toplam maksimum bekleme süresi

		err = backoff.Retry(operation, b)
		if err != nil {
			log.Fatalf("Bucket not ready yet: %v", err)
		}
		collection := bucket.DefaultCollection()
	*/

	bucket := cluster.Bucket("default")
	collection := bucket.DefaultCollection()
	_, err = collection.Upsert(orderID, data, nil)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully recorded..")
}
