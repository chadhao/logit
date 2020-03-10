package response

import (
	"time"

	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserInfoResponse struct {
	Id                 primitive.ObjectID         `json:"id"`
	Phone              string                     `json:"phone"`
	Email              string                     `json:"email"`
	IsEmailVerified    bool                       `json:"isEmailVerified"`
	IsDriver           bool                       `json:"isDriver"`
	RoleIds            []int                      `json:"roleIds"`
	CreatedAt          time.Time                  `json:"createdAt"`
	Driver             *model.Driver              `json:"driver,omitempty"`
	TransportOperators []*model.TransportOperator `json:"transportOperators,omitempty"`
}

func (r *UserInfoResponse) Format(user *model.User, driver *model.Driver, tos []*model.TransportOperator) {
	r.Id = user.Id
	r.Phone = user.Phone
	r.Email = user.Email
	r.IsEmailVerified = user.IsEmailVerified
	r.IsDriver = user.IsDriver
	r.RoleIds = user.RoleIds
	r.CreatedAt = user.CreatedAt
	r.Driver = driver
	r.TransportOperators = tos
}
