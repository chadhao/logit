package model

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mgoClient             *mongo.Client
	db                    *mongo.Database
	suscriptionCollection *mongo.Collection
	recordCollection      *mongo.Collection
	config                map[string]string
	loc                   *time.Location
)

func connect() (err error) {
	uri := config["suscription.db.uri"]
	username := config["suscription.db.username"]
	password := config["suscription.db.password"]
	database := config["suscription.db.database"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mgoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/test?retryWrites=true&w=majority", username, password, uri)
	mgoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mgoURI))
	if err != nil {
		return
	}
	db = mgoClient.Database(database)
	suscriptionCollection = db.Collection("suscription")
	recordCollection = db.Collection("record")
	return
}

// New 创建数据库连接并传入config
func New(c map[string]string) (err error) {
	config = c
	loc, err = time.LoadLocation(config["suscription.time.location"])
	if err != nil {
		return
	}
	err = connect()
	return
}

// Close 关闭
func Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mgoClient.Disconnect(ctx)
}
