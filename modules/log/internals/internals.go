package internals

import (
	"errors"

	"github.com/chadhao/logit/modules/log/model"
	"github.com/chadhao/logit/modules/log/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddLogRequest 创建日志内部请求参数
type AddLogRequest struct {
	Type    string
	Message *string
	FromFun string
	From    *primitive.ObjectID
	Content interface{}
}

// AddLog 创建日志
func AddLog(in *AddLogRequest) error {
	createLogInput := &service.CreateLogInput{
		Message: in.Message,
		FromFun: in.FromFun,
		From:    in.From,
		Content: in.Content,
	}
	var ok = false
	if createLogInput.Type, ok = interface{}(in.Type).(model.Type); !ok {
		return errors.New("type not correct")
	}
	_, err := service.CreateLog(createLogInput)
	return err
}
