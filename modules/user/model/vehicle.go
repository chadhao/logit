package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Vehicle 车辆信息
type Vehicle struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	DriverID     primitive.ObjectID `json:"driverID" bson:"driverID"`
	Registration string             `json:"registration" bson:"registration"`
	IsDiesel     bool               `json:"isDiesel" bson:"isDiesel"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
}

// Create 创建车辆信息
func (v *Vehicle) Create() error {
	_, err := vehicleCollection.InsertOne(context.TODO(), v)
	return err
}

// Delete 删除车辆信息
func (v *Vehicle) Delete() error {
	filter := bson.D{primitive.E{Key: "_id", Value: v.ID}}
	_, err := vehicleCollection.DeleteOne(context.TODO(), filter)
	return err
}

// VehicleExistsOpt 车辆是否存在选项
type VehicleExistsOpt struct {
	DriverID     primitive.ObjectID
	Registration string
}

// IsVehicleExists 车辆是否存在
func IsVehicleExists(opt VehicleExistsOpt) bool {

	conditions := primitive.A{
		bson.D{primitive.E{Key: "driverID", Value: opt.DriverID}},
		bson.D{primitive.E{Key: "registration", Value: opt.Registration}},
	}
	filter := bson.D{primitive.E{Key: "$and", Value: conditions}}
	count, _ := vehicleCollection.CountDocuments(context.TODO(), filter)

	return count > 0
}

// FindVehicle 通过vehicleID获取vehicle信息
func FindVehicle(vehicleID primitive.ObjectID) (*Vehicle, error) {
	vehicle := &Vehicle{}
	err := vehicleCollection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: vehicleID}}).Decode(vehicle)
	return vehicle, err
}

// FindVehiclesOpt 获取vehicles信息选项
type FindVehiclesOpt struct {
	DriverID   primitive.ObjectID
	VehicleIDs []primitive.ObjectID
}

// FindVehicles 获取vehicles信息
func FindVehicles(opt ...FindVehiclesOpt) ([]*Vehicle, error) {
	query := bson.D{}
	if len(opt) == 1 {
		if !opt[0].DriverID.IsZero() {
			query = append(query, primitive.E{Key: "driverID", Value: opt[0].DriverID})
		}
		if len(opt[0].VehicleIDs) > 0 {
			query = append(query, primitive.E{Key: "_id", Value: primitive.E{Key: "$in", Value: opt[0].VehicleIDs}})
		}
	}

	vehicles := []*Vehicle{}
	cursor, err := vehicleCollection.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &vehicles); err != nil {
		return nil, err
	}
	return vehicles, nil
}
