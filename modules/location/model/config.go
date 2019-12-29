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
	dbConfig      map[string]string
	googleConfig  map[string]string
	mgoClient     *mongo.Client
	db            *mongo.Database
	drivingLocCol *mongo.Collection
	mapClient     *maps.Client
)

func dbConnect() (err error) {
	uri := dbConfig["location.db.uri"]
	username := dbConfig["location.db.username"]
	password := dbConfig["location.db.password"]
	database := dbConfig["location.db.database"]

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

// NewDB 创建连接并传入dbConfig
func NewDB(c map[string]string) (err error) {
	dbConfig = c
	return dbConnect()
}

func mapConnect() (err error) {
	mapClient, err = maps.NewClient(maps.WithAPIKey(googleConfig["google.gmap.apikey"]))
	return
}

// NewMap 创建map连接
func NewMap(c map[string]string) (err error) {
	googleConfig = c
	return mapConnect()
}

// Close 关闭
func Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	mgoClient.Disconnect(ctx)
}
