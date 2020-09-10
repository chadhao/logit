package model

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// RedisClient redis client
	RedisClient       *redis.Client
	config            map[string]string
	mongoClient       *mongo.Client
	db                *mongo.Database
	userCollection    *mongo.Collection
	driverCollection  *mongo.Collection
	vehicleCollection *mongo.Collection
	toCollection      *mongo.Collection // transportOperatorCollection
	toICollection     *mongo.Collection // transportOperatorIdentityCollection

)

func connect() (err error) {
	uri := config["user.db.uri"]
	username := config["user.db.username"]
	password := config["user.db.password"]
	database := config["user.db.database"]

	redisAddr := config["user.redis.address"]
	redisPass := config["user.redis.password"]

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mgoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/test?retryWrites=true&w=majority", username, password, uri)
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mgoURI))
	if err != nil {
		return err
	}

	db = mongoClient.Database(database)
	userCollection = db.Collection("user")
	driverCollection = db.Collection("driver")
	vehicleCollection = db.Collection("vehicle")
	toCollection = db.Collection("transportOperator")
	toICollection = db.Collection("transportOperatorIdentity")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       0,
	})

	return nil
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
	ctx, cancel := context.WithCancel(context.Background())
	mongoClient.Disconnect(ctx)
	cancel()
}
