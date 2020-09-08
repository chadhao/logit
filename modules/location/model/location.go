package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (

	// DrivingLoc 司机在行驶过程中的位置信息
	DrivingLoc struct {
		ID        primitive.ObjectID `bson:"_id" json:"id"`
		DriverID  primitive.ObjectID `bson:"driverID" json:"driverID"`
		CreatedAt time.Time          `bson:"createdAt" json:"createdAt" `
		Coors     Coors              `bson:"coors" json:"coors"`
	}
)

// Create 保存司机行驶的位置信息到数据库
func (d *DrivingLoc) Create() error {
	_, err := drivingLocCol.InsertOne(context.TODO(), d)
	return err
}

// GetDrivingLocs 通过driverID和指定时间段返回司机行驶位置信息
func GetDrivingLocs(driverID primitive.ObjectID, from, to time.Time) ([]*DrivingLoc, error) {

	var drivingLocs = []*DrivingLoc{}

	query := bson.M{
		"driverID": driverID,
		"gte":      bson.M{"createdAt": from},
		"lte":      bson.M{"createdAt": to},
	}
	cursor, err := drivingLocCol.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &drivingLocs); err != nil {
		return nil, err
	}

	return drivingLocs, nil
}
