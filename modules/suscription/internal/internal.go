package internal

import "github.com/chadhao/logit/modules/suscription/model"

import "go.mongodb.org/mongo-driver/bson/primitive"

// CreateSuscription 创建用户订阅信息
func CreateSuscription(userID primitive.ObjectID) error {
	s := model.Suscription{
		UserID: userID,
	}
	return s.Create()
}
