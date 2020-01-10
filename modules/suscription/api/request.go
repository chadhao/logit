package api

import "go.mongodb.org/mongo-driver/bson/primitive"

import "github.com/chadhao/logit/modules/suscription/model"

type reqSuscription struct {
	DriverID primitive.ObjectID `json:"driverID" query:"driverID" valid:"required"`
}

func (req *reqSuscription) getSuscription() (*model.Suscription, error) {
	return model.GetSuscription(req.DriverID)
}
