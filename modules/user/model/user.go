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
	filter := bson.D{primitive.E{Key: "_id", Value: u.ID}}
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

	query := bson.D{primitive.E{Key: "$or", Value: conditions}}

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
		query = bson.D{primitive.E{Key: "_id", Value: opt.UserID}}
	case len(opt.Phone) > 0:
		query = bson.D{primitive.E{Key: "phone", Value: opt.Phone}}
	case len(opt.Email) > 0:
		query = bson.D{primitive.E{Key: "email", Value: opt.Email}, primitive.E{Key: "isEmailVerified", Value: true}}
	default:
		return nil, errors.New("No query condition found")
	}

	user := &User{}
	err := userCollection.FindOne(context.TODO(), query).Decode(user)
	return user, err
}

// FilterUserOpt 查找用户选项
type FilterUserOpt struct {
	Phone string
	Email string
}

// FilterUser 按照条件检索user
func FilterUser(opt FilterUserOpt) ([]*User, error) {

	users := []*User{}
	filter := bson.M{}
	if len(opt.Phone) > 0 {
		filter["phone"] = primitive.Regex{Pattern: opt.Phone, Options: "i"}
	}
	if len(opt.Email) > 0 {
		filter["email"] = primitive.Regex{Pattern: opt.Email, Options: "i"}
	}

	projection := bson.D{
		primitive.E{Key: "password", Value: 0},
		primitive.E{Key: "pin", Value: 0},
	}
	cursor, err := userCollection.Find(context.TODO(), filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &users); err != nil {
		return nil, err
	}
	return users, nil
}
