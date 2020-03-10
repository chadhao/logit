package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

type Log struct {
	ID        primitive.ObjectID     `bson:"_id" json:"id" valid:"-"`
	Type      Type                   `json:"type" bson:"type" valid:"required"`
	From      string                 `json:"from" bson:"from" valid:"required"`
	Content   map[string]interface{} `json:"content" bson:"content" valid:"required"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt" valid:"required"`
}
