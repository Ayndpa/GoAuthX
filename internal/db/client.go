package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goauthx/internal/config"
	"sync"
	"time"
)

type MongoConnector struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var (
	clientInstance    *MongoConnector
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
	clientInstance = &MongoConnector{
		Client: client,
		DB:     client.Database(cfg.MongoDB.Database),
	}
}

// GetMongoConnector returns a singleton MongoConnector instance.
func GetMongoConnector() (*MongoConnector, error) {
	mongoOnce.Do(initMongoClient)
	return clientInstance, clientInstanceErr
}