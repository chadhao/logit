package model

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Driver 司机身份信息
type Driver struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	LicenseNumber string             `json:"licenseNumber" bson:"licenseNumber"`
	DateOfBirth   time.Time          `json:"dateOfBirth" bson:"dateOfBirth"`
	Firstnames    string             `json:"firstnames" bson:"firstnames"`
	Surname       string             `json:"surname" bson:"surname"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
}

// Create 创建司机身份
func (d *Driver) Create(user *User) error {
	// 创建司机身份后，更新用户基础信息
	db.Client().UseSession(context.TODO(), func(sessionContext mongo.SessionContext) error {
		// 使用事务
		if err := sessionContext.StartTransaction(); err != nil {
			return err
		}
		if _, err := driverCollection.InsertOne(context.TODO(), d); err != nil {
			return err
		}

		// 用户相关信息更新
		if err := user.Update(); err != nil {
			sessionContext.AbortTransaction(sessionContext)
			return err
		}
		return sessionContext.CommitTransaction(sessionContext)
	})

	return nil
}

// DriverExistsOpt 司机是否存在选项
type DriverExistsOpt struct {
	ID            primitive.ObjectID
	LicenseNumber string
}

// IsDriverExists 判断司机是否存在
func IsDriverExists(opt DriverExistsOpt) bool {
	conditions := bson.D{}
	if !opt.ID.IsZero() {
		conditions = append(conditions, primitive.E{Key: "_id", Value: opt.ID})
	}
	if len(opt.LicenseNumber) > 0 {
		conditions = append(conditions, primitive.E{Key: "licenseNumber", Value: opt.LicenseNumber})
	}
	query := bson.D{primitive.E{Key: "$or", Value: conditions}}

	count, _ := userCollection.CountDocuments(context.TODO(), query)
	return count > 0
}

// FindDriverOpt 查找司机选项
type FindDriverOpt struct {
	ID            primitive.ObjectID
	LicenseNumber string
}

// FindDriver 查找司机
func FindDriver(opt FindDriverOpt) (*Driver, error) {

	query := bson.D{}
	switch {
	case !opt.ID.IsZero():
		query = bson.D{primitive.E{Key: "_id", Value: opt.ID}}
	case len(opt.LicenseNumber) > 0:
		query = bson.D{primitive.E{Key: "licenseNumber", Value: opt.LicenseNumber}}
	default:
		return nil, errors.New("No query condition found")
	}

	driver := &Driver{}
	err := driverCollection.FindOne(context.TODO(), query).Decode(driver)
	return driver, err
}

// FindDrivers 通过IDs查询司机
func FindDrivers(driverIDs []primitive.ObjectID) ([]*Driver, error) {

	drivers := []*Driver{}

	cursor, err := driverCollection.Find(context.TODO(), bson.M{"_id": bson.M{"$in": driverIDs}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &drivers); err != nil {
		return nil, err
	}
	return drivers, nil
}
