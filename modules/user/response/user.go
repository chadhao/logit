package response

import (
	"time"

	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfoResponse struct {
	ID              primitive.ObjectID                `json:"id"`
	Phone           string                            `json:"phone"`
	Email           string                            `json:"email"`
	IsEmailVerified bool                              `json:"isEmailVerified"`
	IsDriver        bool                              `json:"isDriver"`
	RoleIDs         []int                             `json:"roleIDs"`
	CreatedAt       time.Time                         `json:"createdAt"`
	Driver          *model.Driver                     `json:"driver,omitempty"`
	Identities      []model.TransportOperatorIdentity `json:"identities,omitempty"`
}

func (r *UserInfoResponse) Format(user *model.User, driver *model.Driver, identities []model.TransportOperatorIdentity) {
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
	if len(identities) > 0 {
		r.Identities = identities
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
