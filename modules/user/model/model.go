package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		ID              primitive.ObjectID `json:"id,omitempty" bson:"_id"`
		Phone           string             `json:"phone,omitempty" bson:"phone"`
		Email           string             `json:"email,omitempty" bson:"email"`
		IsEmailVerified bool               `json:"isEmailVerified,omitempty" bson:"isEmailVerified"`
		Password        string             `json:"password,omitempty" bson:"password"`
		Pin             string             `json:"pin,omitempty" bson:"pin"`
		IsDriver        bool               `json:"isDriver,omitempty" bson:"isDriver"`
		RoleIDs         []int              `json:"roleIDs,omitempty" bson:"roleIDs"`
		CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	}

	Token struct {
		AccessToken         string             `json:"accessToken,omitempty"`
		AccessTokenExpires  time.Time          `json:"accessTokenExpires,omitempty"`
		RefreshToken        string             `json:"refreshToken,omitempty"`
		RefreshTokenExpires time.Time          `json:"refreshTokenExpires,omitempty"`
		UserID              primitive.ObjectID `json:"userID,omitempty"`
		RoleIDs             []int              `json:"roleIDs,omitempty"`
	}

	Driver struct {
		ID                   primitive.ObjectID   `json:"id" bson:"_id"`
		TransportOperatorIDs []primitive.ObjectID `json:"transportOperatorIDs" bson:"transportOperatorIDs"`
		LicenseNumber        string               `json:"licenseNumber" bson:"licenseNumber"`
		DateOfBirth          time.Time            `json:"dateOfBirth" bson:"dateOfBirth"`
		Firstnames           string               `json:"firstnames" bson:"firstnames"`
		Surname              string               `json:"surname" bson:"surname"`
		CreatedAt            time.Time            `json:"createdAt" bson:"createdAt"`
	}

	Vehicle struct {
		ID           primitive.ObjectID `json:"id" bson:"_id"`
		DriverID     primitive.ObjectID `json:"driverID" bson:"driverID"`
		Registration string             `json:"registration" bson:"registration"`
		IsDiesel     bool               `json:"isDiesel" bson:"isDiesel"`
		CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	}

	TransportOperator struct {
		ID            primitive.ObjectID   `json:"id" bson:"_id"`
		UserIDs       []primitive.ObjectID `json:"userIDs" bson:"userIDs"`
		LicenseNumber string               `json:"licenseNumber" bson:"licenseNumber"`
		Name          string               `json:"name" bson:"name"`
		CreatedAt     time.Time            `json:"createdAt" bson:"createdAt"`
	}
)
