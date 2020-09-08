package service

import (
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/log/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// CreateLogInput 创建日志参数
	CreateLogInput struct {
		Type    model.Type          `json:"type" valid:"required"`
		Message *string             `json:"message,omitempty"`
		FromFun string              `json:"fromFun" valid:"required"`
		From    *primitive.ObjectID `json:"from,omitempty"`
		Content interface{}         `json:"content"`
	}
	// CreateLogOutput 创建日志返回参数
	CreateLogOutput struct {
		*model.Log `json:",inline"`
	}
)

// CreateLog 创建日志
func CreateLog(data *CreateLogInput) (*CreateLogOutput, error) {

	if _, err := valid.ValidateStruct(data); err != nil {
		return nil, err
	}

	l := &model.Log{
		ID:        primitive.NewObjectID(),
		Type:      data.Type,
		Message:   data.Message,
		FromFun:   data.FromFun,
		From:      data.From,
		Content:   data.Content,
		CreatedAt: time.Now(),
	}

	if err := l.Add(); err != nil {
		return nil, err
	}

	return &CreateLogOutput{l}, nil
}

type (
	// QueryLogsInput 查询日志参数
	QueryLogsInput struct {
		model.QueryLogOpt
	}
	// QueryLogsOutput 查询日志返回参数
	QueryLogsOutput struct {
		Logs []*model.Log `json:"logs"`
	}
)

// QueryLogs 查询日志
func QueryLogs(data *QueryLogsInput) (*QueryLogsOutput, error) {
	logs, err := model.QueryLogs(data.QueryLogOpt)
	if err != nil {
		return nil, err
	}
	return &QueryLogsOutput{Logs: logs}, nil
}
