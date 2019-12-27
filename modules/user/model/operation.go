package model

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *User) Create() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if u.Exists() {
		return errors.New("User exists")
	}

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
	if len(u.Id) > 0 && !u.Id.IsZero() {
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

func (u *User) Find() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.D
	if len(u.Id) > 0 && !u.Id.IsZero() {
		filter = bson.D{{"_id", u.Id}}
	} else if len(u.Phone) > 0 {
		filter = bson.D{{"phone", u.Phone}}
	} else if len(u.Email) > 0 {
		filter = bson.D{{"email", u.Email}}
	} else if len(u.DriverId) > 0 && !u.DriverId.IsZero() {
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

func (d *Driver) Find() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.D
	if len(d.Id) > 0 && !d.Id.IsZero() {
		filter = bson.D{{"_id", d.Id}}
	} else if len(d.UserId) > 0 && !d.UserId.IsZero() {
		filter = bson.D{{"userId", d.UserId}}
	} else if len(d.LicenceNumber) > 0 {
		filter = bson.D{{"licenceNumber", d.LicenceNumber}}
	}

	err := db.Collection("driver").FindOne(ctx, filter).Decode(d)

	if err != nil {
		return err
	}

	return nil
}
