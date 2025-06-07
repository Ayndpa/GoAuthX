package db

import (
	"context"
	"sync"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goauthx/pkg/config"
)

var (
	clientInstance    *mongo.Client
	clientInstanceErr error
	mongoOnce         sync.Once
)

func initMongoClient() {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		clientInstanceErr = err
		return
	}
	// Ping to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		clientInstanceErr = err
		return
	}
	clientInstance = client
}

// GetMongoClient returns a singleton MongoDB client instance.
func GetMongoClient() (*mongo.Client, error) {
	mongoOnce.Do(initMongoClient)
	return clientInstance, clientInstanceErr
}