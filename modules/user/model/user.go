package model

import (
	"context"
	"errors"
	"time"

	conf "github.com/chadhao/logit/config"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *User) Create() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if u.Exists() {
		return errors.New("User exists")
	}

	u.Id = primitive.NewObjectID()

	userBson, err := bson.Marshal(u)
	if err != nil {
		return err
	}

	if _, err := db.Collection("user").InsertOne(ctx, userBson); err != nil {
		return err
	}

	return nil
}

func (u *User) Exists() bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conditions := primitive.A{}
	if !u.Id.IsZero() {
		conditions = append(conditions, bson.D{{"_id", u.Id}})
	}
	if len(u.Phone) > 0 {
		conditions = append(conditions, bson.D{{"phone", u.Phone}})
	}
	if len(u.Email) > 0 {
		conditions = append(conditions, bson.D{{"email", u.Email}})
	}

	filter := bson.D{{"$or", conditions}}

	if count, _ := db.Collection("user").CountDocuments(ctx, filter); count > 0 {
		return true
	}

	return false
}

func (u *User) ValidForRegister() bool {
	return len(u.Phone) > 0 && len(u.Password) > 0
}

func (u *User) Find() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.D
	if !u.Id.IsZero() {
		filter = bson.D{{"_id", u.Id}}
	} else if len(u.Phone) > 0 {
		filter = bson.D{{"phone", u.Phone}}
	} else if len(u.Email) > 0 {
		filter = bson.D{{"email", u.Email}}
	} else if !u.DriverId.IsZero() {
		filter = bson.D{{"driverId", u.DriverId}}
	} else {
		return errors.New("No query condition found")
	}

	err := db.Collection("user").FindOne(ctx, filter).Decode(u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) Login() error {
	pass := u.Password

	if err := u.Find(); err != nil {
		return err
	}

	if u.Password != pass {
		return errors.New("Invalid credentials")
	}

	return nil
}

func (u *User) IssueToken(c conf.Config) (*Token, error) {
	now := time.Now().UTC()

	token := &Token{
		AccessTokenExpires:  now.Add(30 * time.Minute),
		RefreshTokenExpires: now.Add(168 * time.Hour),
		UserId:              u.Id,
		RoleIds:             u.RoleIds,
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessTokenClaims := accessToken.Claims.(jwt.MapClaims)
	accessTokenClaims["iss"] = "logit.co.nz"
	accessTokenClaims["exp"] = token.AccessTokenExpires.Unix()
	accessTokenClaims["sub"] = u.Id.Hex()
	accessTokenClaims["roles"] = u.RoleIds
	accessTokenSigningKey, _ := c.Get("system.jwt.access.key")
	if accessTokenSigned, err := accessToken.SignedString([]byte(accessTokenSigningKey)); err != nil {
		return nil, err
	} else {
		token.AccessToken = accessTokenSigned
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["iss"] = "logit.co.nz"
	refreshTokenClaims["exp"] = token.RefreshTokenExpires.Unix()
	refreshTokenClaims["sub"] = u.Id.Hex()
	refreshTokenSigningKey, _ := c.Get("system.jwt.refresh.key")
	if refreshTokenSigned, err := refreshToken.SignedString([]byte(refreshTokenSigningKey)); err != nil {
		return nil, err
	} else {
		token.RefreshToken = refreshTokenSigned
	}

	return token, nil
}
