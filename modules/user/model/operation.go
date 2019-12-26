package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUser(t UserQueryType, v interface{}) (*User, error) {
	user := User{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.D
	switch t {
	case USER_QUERY_BY_ID:
		filter = bson.D{{"_id", v.(primitive.ObjectID)}}
	case USER_QUERY_BY_PHONE:
		filter = bson.D{{"phone", v.(string)}}
	case USER_QUERY_BY_EMAIL:
		filter = bson.D{{"email", v.(string)}}
	case USER_QUERY_BY_LICENCE:
		if driver, err := GetDriver(DRIVER_QUERY_BY_LICENCE, v); err != nil {
			return nil, err
		} else {
			filter = bson.D{{"driverId", driver.Id}}
		}
	}
	err := db.Collection("user").FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetDriver(t DriverQueryType, v interface{}) (*Driver, error) {
	driver := Driver{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.D
	switch t {
	case DRIVER_QUERY_BY_ID:
		filter = bson.D{{"_id", v.(primitive.ObjectID)}}
	case DRIVER_QUERY_BY_LICENCE:
		filter = bson.D{{"licenceNumber", v.(string)}}
	case DRIVER_QUERY_BY_USER_ID:
		filter = bson.D{{"userId", v.(primitive.ObjectID)}}
	}
	err := db.Collection("driver").FindOne(ctx, filter).Decode(&driver)

	if err != nil {
		return nil, err
	}

	return &driver, nil
}
