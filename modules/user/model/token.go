package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token .
type Token struct {
	AccessToken         string             `json:"accessToken,omitempty"`
	AccessTokenExpires  time.Time          `json:"accessTokenExpires,omitempty"`
	RefreshToken        string             `json:"refreshToken,omitempty"`
	RefreshTokenExpires time.Time          `json:"refreshTokenExpires,omitempty"`
	UserID              primitive.ObjectID `json:"userID,omitempty"`
	RoleIDs             []int              `json:"roleIDs,omitempty"`
}
