package model

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"googlemaps.github.io/maps"
)

var (
	mgoClient     *mongo.Client
	db            *mongo.Database
	drivingLocCol *mongo.Collection
	config        map[string]string
	mapClient     *maps.Client
)

func dbConnect() (err error) {
	uri := config["location.db.uri"]
	username := config["location.db.username"]
	password := config["location.db.password"]
	database := config["location.db.database"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mgoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/test?retryWrites=true&w=majority", username, password, uri)
	mgoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mgoURI))
	if err != nil {
		return
	}
	db = mgoClient.Database(database)
	drivingLocCol = db.Collection("driving_location")
	return
}

// New 创建连接并传入config
func New(c map[string]string) (err error) {
	config = c

	if err = dbConnect(); err != nil {
		return err
	}

	mapClient, err = maps.NewClient(maps.WithAPIKey(config["location.gmap.apikey"]))
	return
}

// Close 关闭
func Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mgoClient.Disconnect(ctx)
}
