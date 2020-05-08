package model

import (
	"time"

	valid "github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type (
	// Type 记录类型
	Type string
)

const (
	// Error 错误记录类型
	Error Type = "error"
	// Info 信息记录类型
	Info Type = "info"
	// Modification 变更记录类型
	Modification Type = "modification"
)

// Log .
type Log struct {
	ID        primitive.ObjectID  `bson:"_id" json:"id" valid:"-"`
	Type      Type                `json:"type" bson:"type" valid:"required"`
	FromFun   string              `json:"fromFun" bson:"fromFun" valid:"required"`
	From      *primitive.ObjectID `json:"from,omitempty" bson:"from,omitempty"`
	Content   interface{}         `json:"content" bson:"content" valid:"required"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt" valid:"required"`
}

// Add 数据库添加记录
func (l *Log) Add() error {
	if _, err := valid.ValidateStruct(l); err != nil {
		return err
	}
	if _, err := logCollection.InsertOne(context.TODO(), l); err != nil {
		return err
	}
	return nil
}

// QueryLog .
type QueryLog struct {
	Type    *Type     `json:"type,omitempty"`
	FromFun *string   `json:"fromFun,omitempty"`
	From    time.Time `json:"from" valid:"required"`
	To      time.Time `json:"to" valid:"required"`
}

// Find .
func (q *QueryLog) Find() ([]Log, error) {
	if _, err := valid.ValidateStruct(q); err != nil {
		return nil, err
	}
	logs := []Log{}
	filter := bson.M{
		"createdAt": bson.M{"$gte": q.From, "$lte": q.To},
	}
	if q.Type != nil {
		filter["type"] = *q.Type
	}
	if q.FromFun != nil {
		filter["fromFun"] = *q.FromFun
	}
	cursor, err := logCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
