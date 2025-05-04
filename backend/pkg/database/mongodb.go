package database

import (
	"context"
	"log"
	"time"

	"github.com/qiniu/qmgo"
)

func NewMongoClient(uri string, db string) (*qmgo.Database, error) {
	log.Println("Initializing MongoDB connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: uri})
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v\n", err)
		return nil, err
	}

	log.Println("Successfully connected to MongoDB. Packaging database...")
	database := client.Database(db)
	log.Printf("Database '%s' is ready for use.\n", db)

	return database, nil
}
