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
