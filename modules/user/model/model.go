package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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
		LicenceNumber        string               `json:"licenceNumber" bson:"licenceNumber"`
		DateOfBirth          time.Time            `json:"dateOfBirth" bson:"dateOfBirth"`
		Firstnames           string               `json:"firstnames" bson:"firstnames"`
		Surname              string               `json:"surname" bson:"surname"`
		TransportOperatorIds []primitive.ObjectID `json:"transportOperatorIds" bson:"transportOperatorIds"`
	}

	TransportOperator struct {
		Id            primitive.ObjectID   `json:"id" bson:"_id"`
		LicenceNumber string               `json:"licenceNumber" bson:"licenceNumber"`
		Name          string               `json:"name" bson:"name"`
		UserIds       []primitive.ObjectID `json:"userIds" bson:"userIds"`
	}
)
