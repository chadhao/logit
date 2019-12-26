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
		Id       primitive.ObjectID   `json:"id" bson:"_id"`
		Phone    string               `json:"phone" bson:"phone"`
		Email    string               `json:"email" bson:"email"`
		Password string               `json:"password" bson:"password"`
		Pin      string               `json:"pin" bson:"pin"`
		DriverId primitive.ObjectID   `json:"driverId" bson:"driverId"`
		RoleIds  []primitive.ObjectID `json:"roleIds" bson:"roleIds"`
	}

	Role struct {
		Id   primitive.ObjectID `json:"id" bson:"_id"`
		Name string             `json:"name" bson:"name"`
	}

	Driver struct {
		Id                   primitive.ObjectID   `json:"id" bson:"_id"`
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
