package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Transaction 交易结构体
type Transaction struct {
	ID       primitive.ObjectID `bson:"_id" json:"id" valid:"required"`
	DriverID primitive.ObjectID `bson:"driverID" json:"driverID" valid:"required"`
	RecordID primitive.ObjectID `bson:"recordID" json:"recordID" valid:"required"`
}
