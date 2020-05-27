package request

import (
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/user/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	VehicleCreateRequest struct {
		DriverID     primitive.ObjectID `json:"driverID" valid:"required"`
		Registration string             `json:"registration" valid:"stringlength(5|8)"`
		IsDiesel     bool               `json:"isDiesel" valid:"-"`
	}
)

func (r *VehicleCreateRequest) Create() (*model.Vehicle, error) {
	if _, err := valid.ValidateStruct(r); err != nil {
		return nil, err
	}

	vehicle := &model.Vehicle{
		ID:           primitive.NewObjectID(),
		DriverID:     r.DriverID,
		Registration: r.Registration,
		IsDiesel:     r.IsDiesel,
		CreatedAt:    time.Now(),
	}

	return vehicle, vehicle.Create()
}
