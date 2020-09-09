package model

import (
	"time"

	"github.com/chadhao/logit/modules/user/constant"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	TransportOperator struct {
		ID            primitive.ObjectID `json:"id" bson:"_id"`
		LicenseNumber string             `json:"licenseNumber" bson:"licenseNumber"`
		Name          string             `json:"name" bson:"name"`
		IsVerified    bool               `json:"isVerified" bson:"isVerified"`
		IsCompany     bool               `json:"isCompany" bson:"isCompany"`
		CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	}

	TransportOperatorIdentity struct {
		ID                  primitive.ObjectID `json:"id" bson:"_id"`
		UserID              primitive.ObjectID `json:"userID" bson:"userID"`
		TransportOperatorID primitive.ObjectID `json:"transportOperatorID" bson:"transportOperatorID"`
		Identity            TOIdentity         `json:"identity" bson:"identity"`
		Contact             *string            `json:"contact" bson:"contact"`
		CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
	}

	TransportOperatorIdentityDetail struct {
		ID                  primitive.ObjectID `json:"id" bson:"_id"`
		UserID              primitive.ObjectID `json:"userID" bson:"userID"`
		TransportOperatorID primitive.ObjectID `json:"transportOperatorID" bson:"transportOperatorID"`
		Identity            TOIdentity         `json:"identity" bson:"identity"`
		Contact             *string            `json:"contact" bson:"contact"`
		CreatedAt           time.Time          `json:"createdAt" bson:"createdAt"`
		TransportOperator   *TransportOperator `json:"transportOperator" bson:"transportOperator"`
	}
)

type (
	// TOIdentity transport operator identity
	TOIdentity string
)

const (
	// TO_ADMIN admin
	TO_ADMIN TOIdentity = "to_admin"
	// TO_SUPER super
	TO_SUPER TOIdentity = "to_super"
	// TO_DRIVER driver
	TO_DRIVER TOIdentity = "to_driver"
)

func (t TOIdentity) GetRole() int {
	identity := -1
	switch {
	case t == TO_SUPER:
		identity = constant.ROLE_TO_SUPER
	case t == TO_ADMIN:
		identity = constant.ROLE_TO_ADMIN
	case t == TO_DRIVER:
		identity = constant.ROLE_DRIVER
	}
	return identity
}
