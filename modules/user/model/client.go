package model

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	config      map[string]string
	mongoClient *mongo.Client
	redisClient *redis.Client
	db          *mongo.Database
)

func connect() error {
	uri := config["user.db.uri"]
	username := config["user.db.username"]
	password := config["user.db.password"]
	database := config["user.db.database"]

	redisAddr := config["user.redis.address"]
	redisPass := config["user.redis.password"]

	var err error
	ctx, cancel := context.WithCancel(context.Background())
	connectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s/test?retryWrites=true&w=majority", username, password, uri)
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	cancel()

	if err != nil {
		return err
	}

	db = mongoClient.Database(database)
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})

	return nil
}

func New(c map[string]string) error {
	config = c

	if err := connect(); err != nil {
		return err
	}

	return nil
}

func Close() {
	ctx, cancel := context.WithCancel(context.Background())
	mongoClient.Disconnect(ctx)
	cancel()
}
