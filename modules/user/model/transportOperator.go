package model

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (t *TransportOperator) Create() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if t.Exists() {
		return errors.New("Transport operator exists")
	}

	toBson, err := bson.Marshal(t)
	if err != nil {
		return err
	}

	if _, err := db.Collection("transportOperator").InsertOne(ctx, toBson); err != nil {
		return err
	}

	return nil
}

func (t *TransportOperator) Exists() bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conditions := primitive.A{}
	if !t.Id.IsZero() {
		conditions = append(conditions, bson.D{{"_id", t.Id}})
	}
	if len(t.LicenseNumber) > 0 {
		conditions = append(conditions, bson.D{{"licenseNumber", t.LicenseNumber}})
	}

	filter := bson.D{{"$or", conditions}}

	if count, _ := db.Collection("transportOperator").CountDocuments(ctx, filter); count > 0 {
		return true
	}

	return false
}

func (t *TransportOperator) Find() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.D
	if !t.Id.IsZero() {
		filter = bson.D{{"_id", t.Id}}
	} else if len(t.LicenseNumber) > 0 {
		filter = bson.D{{"licenseNumber", t.LicenseNumber}}
	}

	err := db.Collection("transportOperator").FindOne(ctx, filter).Decode(t)

	if err != nil {
		return err
	}

	return nil
}
