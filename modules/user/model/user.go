package model

import (
	"context"
	"errors"
	"time"

	conf "github.com/chadhao/logit/config"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (u *User) Create() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if u.Exists() {
		return errors.New("User exists")
	}

	u.ID = primitive.NewObjectID()

	userBson, err := bson.Marshal(u)
	if err != nil {
		return err
	}

	if _, err := db.Collection("user").InsertOne(ctx, userBson); err != nil {
		return err
	}

	return nil
}

func (u *User) Update() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.D{{"_id", u.ID}}
	userBson, err := bson.Marshal(u)
	if err != nil {
		return err
	}

	result, err := db.Collection("user").ReplaceOne(ctx, filter, userBson)
	if err != nil {
		return nil
	}
	if result.MatchedCount != 1 {
		return errors.New("User not updated")
	}

	return nil
}

func (u *User) Exists() bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conditions := primitive.A{}
	if !u.ID.IsZero() {
		conditions = append(conditions, bson.D{{"_id", u.ID}})
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

func (u *User) Find() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.D

	if !u.ID.IsZero() {
		filter = bson.D{{"_id", u.ID}}
	} else if len(u.Phone) > 0 {
		filter = bson.D{{"phone", u.Phone}}
	} else if len(u.Email) > 0 {
		filter = bson.D{{"email", u.Email}, {"isEmailVerified", true}}
	} else {
		return errors.New("No query condition found")
	}
	err := db.Collection("user").FindOne(ctx, filter).Decode(u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordLogin() error {
	pass := u.Password

	if err := u.Find(); err != nil {
		return errors.New("user not found")
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
		UserID:              u.ID,
		RoleIDs:             u.RoleIDs,
	}

	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessTokenClaims := accessToken.Claims.(jwt.MapClaims)
	accessTokenClaims["iss"] = "logit.co.nz"
	accessTokenClaims["exp"] = token.AccessTokenExpires.Unix()
	accessTokenClaims["sub"] = u.ID.Hex()
	accessTokenClaims["roles"] = u.RoleIDs
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
	refreshTokenClaims["sub"] = u.ID.Hex()
	refreshTokenSigningKey, _ := c.Get("system.jwt.refresh.key")
	if refreshTokenSigned, err := refreshToken.SignedString([]byte(refreshTokenSigningKey)); err != nil {
		return nil, err
	} else {
		token.RefreshToken = refreshTokenSigned
	}

	return token, nil
}

// Filter 按照条件检索user
func (u *User) Filter() ([]User, error) {

	users := []User{}

	filter := bson.M{}
	if len(u.Phone) > 0 {
		filter["phone"] = primitive.Regex{Pattern: u.Phone, Options: "i"}
	}
	if len(u.Email) > 0 {
		filter["email"] = primitive.Regex{Pattern: u.Email, Options: "i"}
	}

	projection := bson.D{
		{"password", 0},
		{"pin", 0},
	}
	cursor, err := db.Collection("user").Find(context.TODO(), filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}
