package record

import (
	"context"
	"time"

	"github.com/chadhao/logit/config"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	// mgoClient        *mongo.Client
	recordDB         *mongo.Database
	recordCollection *mongo.Collection
)

// InitModule 模块初始化
func InitModule(e *echo.Echo, c *config.Config) {
	// load config
	// add routes
	// other initialization code
}

func initMongoDB(mgoURI, mgoDBName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mgoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mgoURI))
	if err != nil {
		return err
	}
	recordDB = mgoClient.Database(mgoDBName)
	recordCollection = recordDB.Collection("record")
	return nil
}
