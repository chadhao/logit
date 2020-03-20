package response

import (
	"time"

	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfoResponse struct {
	ID                 primitive.ObjectID `json:"id"`
	Phone              string             `json:"phone"`
	Email              string             `json:"email"`
	IsEmailVerified    bool               `json:"isEmailVerified"`
	IsDriver           bool               `json:"isDriver"`
	RoleIDs            []int              `json:"roleIDs"`
	CreatedAt          time.Time          `json:"createdAt"`
	Driver             *model.Driver      `json:"driver,omitempty"`
	TransportOperators struct {
		Drivers []model.TransportOperator `json:"drivers,omitempty"`
		Supers  []model.TransportOperator `json:"supers,omitempty"`
		Admins  []model.TransportOperator `json:"admins,omitempty"`
	} `json:"transportOperators,omitempty"`
}

func (r *UserInfoResponse) AddToInfo(drivers []model.TransportOperator, supers []model.TransportOperator, admins []model.TransportOperator) {
	if len(drivers) > 0 {
		r.TransportOperators.Drivers = drivers
	}
	if len(supers) > 0 {
		r.TransportOperators.Supers = supers
	}
	if len(admins) > 0 {
		r.TransportOperators.Admins = admins
	}
}

func (r *UserInfoResponse) Format(user *model.User, driver *model.Driver) {
	r.ID = user.ID
	r.Phone = user.Phone
	r.Email = user.Email
	r.IsEmailVerified = user.IsEmailVerified
	r.IsDriver = user.IsDriver
	r.RoleIDs = user.RoleIDs
	r.CreatedAt = user.CreatedAt
	if !driver.ID.IsZero() {
		r.Driver = driver
	}
}

type TransportOperatorInfoResponse struct {
	ID            primitive.ObjectID `json:"id"`
	LicenseNumber string             `json:"licenseNumber"`
	Name          string             `json:"name"`
	IsVerified    bool               `json:"isVerified"`
}

func TransportOperatorInfoFormat(tos []model.TransportOperator) []TransportOperatorInfoResponse {
	r := []TransportOperatorInfoResponse{}
	for _, v := range tos {
		t := TransportOperatorInfoResponse{
			ID:            v.ID,
			LicenseNumber: v.LicenseNumber,
			Name:          v.Name,
			IsVerified:    v.IsVerified,
		}
		r = append(r, t)
	}
	return r
}
