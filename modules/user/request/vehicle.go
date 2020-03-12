package request

import (
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/user/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	VehicleCreateRequest struct {
		DriverId     primitive.ObjectID `json:"driverId" valid:"-"`
		Registration string             `json:"registration" valid:"numeric,stringlength(5|9)"`
		IsDiesel     bool               `json:"isDiesel" valid:"required"`
	}
)

func (r *VehicleCreateRequest) Create() (*model.Vehicle, error) {
	if _, err := valid.ValidateStruct(r); err != nil {
		return nil, err
	}

	vehicle := &model.Vehicle{
		Id:           primitive.NewObjectID(),
		DriverId:     r.DriverId,
		Registration: r.Registration,
		IsDiesel:     r.IsDiesel,
		CreatedAt:    time.Now(),
	}

	return vehicle, vehicle.Create()
}
