package model

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mgoClient        *mongo.Client
	db               *mongo.Database
	recordCollection *mongo.Collection
	noteCollection   *mongo.Collection
	config           map[string]string
)

func connect() (err error) {
	uri := config["record.db.uri"]
	username := config["record.db.username"]
	password := config["record.db.password"]
	database := config["record.db.database"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mgoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/test?retryWrites=true&w=majority", username, password, uri)
	mgoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mgoURI))
	if err != nil {
		return
	}
	db = mgoClient.Database(database)
	recordCollection = db.Collection("record")
	noteCollection = db.Collection("note")
	return
}

// New 创建数据库连接并传入config
func New(c map[string]string) error {
	config = c

	if err := connect(); err != nil {
		return err
	}

	return nil
}

// Close 关闭
func Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mgoClient.Disconnect(ctx)
}
