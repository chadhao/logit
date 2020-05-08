package internals

import (
	"time"

	"github.com/chadhao/logit/modules/log/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReqAddLog .
type ReqAddLog struct {
	Type    string              `json:"type" bson:"type" valid:"required"`
	FromFun string              `json:"fromFun" bson:"fromFun" valid:"required"`
	From    *primitive.ObjectID `json:"from,omitempty" bson:"from,omitempty"`
	Content interface{}         `json:"content" bson:"content" valid:"required"`
}

// AddLog 添加log
func AddLog(r *ReqAddLog) error {
	log := &model.Log{
		ID:        primitive.NewObjectID(),
		Type:      interface{}(r.Type).(model.Type),
		FromFun:   r.FromFun,
		From:      r.From,
		Content:   r.Content,
		CreatedAt: time.Now(),
	}
	return log.Add()
}
