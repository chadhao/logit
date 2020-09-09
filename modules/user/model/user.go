package model

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User 基础用户信息
type User struct {
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

// Create 创建基础用户
func (u *User) Create() error {
	_, err := userCollection.InsertOne(context.TODO(), u)
	return err
}

// Update 更新基础用户信息
func (u *User) Update() error {
	filter := bson.D{{"_id", u.ID}}
	if result, _ := userCollection.ReplaceOne(context.TODO(), filter, u); result.MatchedCount != 1 {
		return errors.New("User not updated")
	}
	return nil
}

// UserExistsOpt 判断用户是否存在选项
type UserExistsOpt struct {
	ID    primitive.ObjectID
	Phone string
	Email string
}

// IsUserExists 判断用户是否存在
func IsUserExists(opt UserExistsOpt) bool {
	conditions := bson.D{}
	if !opt.ID.IsZero() {
		conditions = append(conditions, primitive.E{Key: "_id", Value: opt.ID})
	}
	if len(opt.Phone) > 0 {
		conditions = append(conditions, primitive.E{Key: "phone", Value: opt.Phone})
	}
	if len(opt.Email) > 0 {
		conditions = append(conditions, primitive.E{Key: "email", Value: opt.Email})
	}

	query := bson.D{{"$or", conditions}}

	count, _ := userCollection.CountDocuments(context.TODO(), query)
	return count > 0
}

// FindUserOpt 查找用户选项
type FindUserOpt struct {
	UserID primitive.ObjectID
	Phone  string
	Email  string
}

// FindUser 查找用户
func FindUser(opt FindUserOpt) (*User, error) {
	query := bson.D{}
	switch {
	case !opt.UserID.IsZero():
		query = bson.D{{"_id", opt.UserID}}
	case len(opt.Phone) > 0:
		query = bson.D{{"phone", opt.Phone}}
	case len(opt.Email) > 0:
		query = bson.D{{"email", opt.Email}, {"isEmailVerified", true}}
	default:
		return nil, errors.New("No query condition found")
	}

	user := &User{}
	err := userCollection.FindOne(context.TODO(), query).Decode(user)
	return user, err
}

// func (u *User) Find() error {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	var filter bson.D

// 	if !u.ID.IsZero() {
// 		filter = bson.D{{"_id", u.ID}}
// 	} else if len(u.Phone) > 0 {
// 		filter = bson.D{{"phone", u.Phone}}
// 	} else if len(u.Email) > 0 {
// 		filter = bson.D{{"email", u.Email}, {"isEmailVerified", true}}
// 	} else {
// 		return errors.New("No query condition found")
// 	}
// 	err := db.Collection("user").FindOne(ctx, filter).Decode(u)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (u *User) PasswordLogin() error {
// 	pass := u.Password

// 	if err := u.Find(); err != nil {
// 		return errors.New("user not found")
// 	}

// 	if u.Password != pass {
// 		return errors.New("Invalid credentials")
// 	}

// 	return nil
// }

// IssueToken 给该用户生成token
// func (u *User) IssueToken(c conf.Config) (*Token, error) {
// 	now := time.Now().UTC()

// 	token := &Token{
// 		AccessTokenExpires:  now.Add(30 * time.Minute),
// 		RefreshTokenExpires: now.Add(168 * time.Hour),
// 		UserID:              u.ID,
// 		RoleIDs:             u.RoleIDs,
// 	}

// 	accessToken := jwt.New(jwt.SigningMethodHS256)
// 	accessTokenClaims := accessToken.Claims.(jwt.MapClaims)
// 	accessTokenClaims["iss"] = "logit.co.nz"
// 	accessTokenClaims["exp"] = token.AccessTokenExpires.Unix()
// 	accessTokenClaims["sub"] = u.ID.Hex()
// 	accessTokenClaims["roles"] = u.RoleIDs
// 	accessTokenSigningKey, _ := c.Get("system.jwt.access.key")
// 	if accessTokenSigned, err := accessToken.SignedString([]byte(accessTokenSigningKey)); err != nil {
// 		return nil, err
// 	} else {
// 		token.AccessToken = accessTokenSigned
// 	}

// 	refreshToken := jwt.New(jwt.SigningMethodHS256)
// 	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
// 	refreshTokenClaims["iss"] = "logit.co.nz"
// 	refreshTokenClaims["exp"] = token.RefreshTokenExpires.Unix()
// 	refreshTokenClaims["sub"] = u.ID.Hex()
// 	refreshTokenSigningKey, _ := c.Get("system.jwt.refresh.key")
// 	if refreshTokenSigned, err := refreshToken.SignedString([]byte(refreshTokenSigningKey)); err != nil {
// 		return nil, err
// 	} else {
// 		token.RefreshToken = refreshTokenSigned
// 	}

// 	return token, nil
// }

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
