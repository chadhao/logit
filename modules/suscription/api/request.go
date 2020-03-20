package api

import (
	"github.com/chadhao/logit/modules/suscription/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type reqSuscription struct {
	DriverID string `json:"driverID" query:"driverID" valid:"required"`
}

func (req *reqSuscription) getSuscription() (*model.Suscription, error) {
	driverID, err := primitive.ObjectIDFromHex(req.DriverID)
	if err != nil {
		return nil, err
	}
	return model.GetSuscription(driverID)
}
