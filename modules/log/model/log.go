package model

import (
	"time"

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
	ID        primitive.ObjectID  `json:"id" bson:"_id"`
	Type      Type                `json:"type" bson:"type"`
	Message   *string             `json:"message,omitempty" bson:"message,omitempty"`
	FromFun   string              `json:"fromFun" bson:"fromFun"`
	From      *primitive.ObjectID `json:"from,omitempty" bson:"from,omitempty"`
	Content   interface{}         `json:"content" bson:"content"`
	CreatedAt time.Time           `json:"createdAt" bson:"createdAt"`
}

// Add 数据库添加记录
func (l *Log) Add() error {
	_, err := logCollection.InsertOne(context.TODO(), l)
	return err
}

// QueryLogOpt 查询日志选项
type QueryLogOpt struct {
	Type    *Type     `json:"type,omitempty"`
	FromFun *string   `json:"fromFun,omitempty"`
	From    time.Time `json:"from"`
	To      time.Time `json:"to"`
}

func (q *QueryLogOpt) query() bson.D {
	query := bson.D{}
	if q.Type != nil {
		query = append(query, primitive.E{Key: "type", Value: *q.Type})
	}
	if q.FromFun != nil {
		query = append(query, primitive.E{Key: "fromFun", Value: *q.FromFun})
	}
	if !q.From.IsZero() {
		query = append(query, primitive.E{Key: "createdAt", Value: primitive.E{Key: "$gte", Value: q.From}})
	}
	if !q.To.IsZero() {
		query = append(query, primitive.E{Key: "createdAt", Value: primitive.E{Key: "$lte", Value: q.To}})
	}
	return query
}

// QueryLogs 查询日志
func QueryLogs(opt QueryLogOpt) ([]*Log, error) {
	query := opt.query()
	if len(query) == 0 {
		query = nil
	}

	logs := []*Log{}
	cursor, err := logCollection.Find(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
