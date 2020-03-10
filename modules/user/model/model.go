package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		Id       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
		Phone    string             `json:"phone,omitempty" bson:"phone"`
		Email    string             `json:"email,omitempty" bson:"email"`
		Password string             `json:"password,omitempty" bson:"password"`
		Pin      string             `json:"pin,omitempty" bson:"pin"`
		IsDriver bool               `json:"isDriver,omitempty" bson:"isDriver"`
		RoleIds  []int              `json:"roleIds,omitempty" bson:"roleIds"`
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
		TransportOperatorIds []primitive.ObjectID `json:"transportOperatorIds" bson:"transportOperatorIds"`
		LicenceNumber        string               `json:"licenceNumber" bson:"licenceNumber"`
		DateOfBirth          time.Time            `json:"dateOfBirth" bson:"dateOfBirth"`
		Firstnames           string               `json:"firstnames" bson:"firstnames"`
		Surname              string               `json:"surname" bson:"surname"`
	}

	Vehicle struct {
		Id           primitive.ObjectID `json:"id" bson:"_id"`
		DriverId     primitive.ObjectID `json:"driverId" bson:"driverId"`
		Registration string             `json:"registration" bson:"registration"`
		IsDiesel     bool               `json:"isDiesel" bson:"isDiesel"`
	}

	TransportOperator struct {
		Id            primitive.ObjectID   `json:"id" bson:"_id"`
		UserIds       []primitive.ObjectID `json:"userIds" bson:"userIds"`
		LicenceNumber string               `json:"licenceNumber" bson:"licenceNumber"`
		Name          string               `json:"name" bson:"name"`
	}
)
