package model

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (v *Vehicle) Create() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if v.Exists() {
		return errors.New("Vehicle exists")
	}

	v.Id = primitive.NewObjectID()

	vehicleBson, err := bson.Marshal(v)
	if err != nil {
		return err
	}

	if _, err := db.Collection("vehicle").InsertOne(ctx, vehicleBson); err != nil {
		return err
	}

	return nil
}

func (v *Vehicle) Delete() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.D{{"_id", v.Id}}

	if _, err := db.Collection("vehicle").DeleteOne(ctx, filter); err != nil {
		return nil
	}

	return nil
}

func (v *Vehicle) Exists() bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conditions := primitive.A{
		bson.D{{"driverId", v.DriverId}},
		bson.D{{"registration", v.Registration}},
	}

	filter := bson.D{{"$and", conditions}}

	if count, _ := db.Collection("vehicle").CountDocuments(ctx, filter); count > 0 {
		return true
	}

	return false
}

func (v *Vehicle) Find() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.D{{"_id", v.Id}}

	err := db.Collection("vehicle").FindOne(ctx, filter).Decode(v)

	if err != nil {
		return err
	}

	return nil
}

func (v *Vehicle) FindByDriverId() ([]Vehicle, error) {
	vehicles := []Vehicle{}
	filter := bson.M{
		"driverId": v.DriverId,
	}

	cursor, err := db.Collection("vehicle").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &vehicles); err != nil {
		return nil, err
	}
	return vehicles, nil
}
