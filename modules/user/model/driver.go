package model

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *Driver) Create() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if d.Exists() {
		return errors.New("Driver exists")
	}

	driverBson, err := bson.Marshal(d)
	if err != nil {
		return err
	}

	if _, err := db.Collection("driver").InsertOne(ctx, driverBson); err != nil {
		return err
	}

	return nil
}

func (d *Driver) Exists() bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conditions := primitive.A{}
	if !d.Id.IsZero() {
		conditions = append(conditions, bson.D{{"_id", d.Id}})
	}
	if len(d.LicenseNumber) > 0 {
		conditions = append(conditions, bson.D{{"licenseNumber", d.LicenseNumber}})
	}

	filter := bson.D{{"$or", conditions}}

	if count, _ := db.Collection("driver").CountDocuments(ctx, filter); count > 0 {
		return true
	}

	return false
}

func (d *Driver) Find() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.D
	if !d.Id.IsZero() {
		filter = bson.D{{"_id", d.Id}}
	} else if len(d.LicenseNumber) > 0 {
		filter = bson.D{{"licenseNumber", d.LicenseNumber}}
	}

	err := db.Collection("driver").FindOne(ctx, filter).Decode(d)

	if err != nil {
		return err
	}

	return nil
}
