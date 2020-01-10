package internal

import "github.com/chadhao/logit/modules/suscription/model"

import "go.mongodb.org/mongo-driver/bson/primitive"

import "time"

// CreateSuscription 创建用户订阅信息
func CreateSuscription(driverID primitive.ObjectID, renew bool) error {
	s := &model.Suscription{
		DriverID:  driverID,
		Renew:     renew,
		CreatedAt: time.Now(),
	}
	return s.Add()
}

// RespSuscirption 系统内部订阅状态返回
type RespSuscirption struct {
	model.Suscription `json:",inline"`
	IsExpired         bool `json:"isExpired"`
}

// GetSuscription 获取订阅状态
func GetSuscription(driverID primitive.ObjectID) (*RespSuscirption, error) {
	s, err := model.GetSuscription(driverID)
	if err != nil {
		return nil, err
	}
	resp := &RespSuscirption{
		Suscription: *s,
		IsExpired:   s.IsExpired(),
	}
	return resp, nil
}
