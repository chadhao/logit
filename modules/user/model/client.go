package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config map[string]string
var mongoClient *mongo.Client
var db *mongo.Database

func connect() error {
	uri := config["user.db.uri"]
	username := config["user.db.username"]
	password := config["user.db.password"]
	database := config["user.db.database"]

	var err error
	ctx, cancel := context.WithCancel(context.Background())
	connectionString := fmt.Sprintf("mongodb+srv://%s:%s@%s/test?retryWrites=true&w=majority", username, password, uri)
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	cancel()

	if err != nil {
		return err
	}

	db = mongoClient.Database(database)

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
