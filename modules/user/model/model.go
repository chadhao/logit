package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type (
	UserQueryType   int
	DriverQueryType int
)

const (
	USER_QUERY_BY_ID      UserQueryType = 0
	USER_QUERY_BY_PHONE   UserQueryType = 1
	USER_QUERY_BY_EMAIL   UserQueryType = 2
	USER_QUERY_BY_LICENCE UserQueryType = 3

	DRIVER_QUERY_BY_ID      DriverQueryType = 0
	DRIVER_QUERY_BY_LICENCE DriverQueryType = 1
	DRIVER_QUERY_BY_USER_ID DriverQueryType = 2
)

type (
	User struct {
		Id              primitive.ObjectID `json:"id,omitempty" bson:"_id"`
		Phone           string             `json:"phone,omitempty" bson:"phone"`
		Email           string             `json:"email,omitempty" bson:"email"`
		Password        string             `json:"password,omitempty" bson:"password"`
		ConfirmPassword string             `json:"confirmPassword,omitempty" bson:"-"`
		Pin             string             `json:"pin,omitempty" bson:"pin"`
		DriverId        primitive.ObjectID `json:"driverId,omitempty" bson:"driverId"`
		RoleIds         []int              `json:"roleIds,omitempty" bson:"roleIds"`
	}

	Token struct {
		AccessToken         string             `json:"accessToken,omitempty"`
		AccessTokenExpires  time.Time          `json:"accessTokenExpires,omitempty"`
		RefreshToken        string             `json:"refreshToken,omitempty"`
		RefreshTokenExpires time.Time          `json:"refreshTokenExpires,omitempty"`
		UserId              primitive.ObjectID `json:"userId,omitempty"`
		RoleIds             []int              `json:"roleIds,omitempty"`
	}

	Driver struct {
		Id                   primitive.ObjectID   `json:"id" bson:"_id"`
		UserId               primitive.ObjectID   `json:"userId" bson:"userId"`
		TransportOperatorIds []primitive.ObjectID `json:"transportOperatorIds" bson:"transportOperatorIds"`
		LicenceNumber        string               `json:"licenceNumber" bson:"licenceNumber"`
		DateOfBirth          time.Time            `json:"dateOfBirth" bson:"dateOfBirth"`
		Firstnames           string               `json:"firstnames" bson:"firstnames"`
		Surname              string               `json:"surname" bson:"surname"`
	}

	TransportOperator struct {
		Id            primitive.ObjectID   `json:"id" bson:"_id"`
		UserIds       []primitive.ObjectID `json:"userIds" bson:"userIds"`
		LicenceNumber string               `json:"licenceNumber" bson:"licenceNumber"`
		Name          string               `json:"name" bson:"name"`
	}
)
