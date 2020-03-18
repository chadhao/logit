package response

import (
	"time"

	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfoResponse struct {
	ID                 primitive.ObjectID         `json:"id"`
	Phone              string                     `json:"phone"`
	Email              string                     `json:"email"`
	IsEmailVerified    bool                       `json:"isEmailVerified"`
	IsDriver           bool                       `json:"isDriver"`
	RoleIDs            []int                      `json:"roleIDs"`
	CreatedAt          time.Time                  `json:"createdAt"`
	Driver             *model.Driver              `json:"driver,omitempty"`
	TransportOperators []*model.TransportOperator `json:"transportOperators,omitempty"`
}

func (r *UserInfoResponse) Format(user *model.User, driver *model.Driver, tos []*model.TransportOperator) {
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
	if len(tos) > 0 {
		r.TransportOperators = tos
	}
}
