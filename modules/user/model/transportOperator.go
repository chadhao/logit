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
	if !t.ID.IsZero() {
		conditions = append(conditions, bson.D{{"_id", t.ID}})
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
	if !t.ID.IsZero() {
		filter = bson.D{{"_id", t.ID}}
	} else if len(t.LicenseNumber) > 0 {
		filter = bson.D{{"licenseNumber", t.LicenseNumber}}
	}

	err := db.Collection("transportOperator").FindOne(ctx, filter).Decode(t)

	if err != nil {
		return err
	}

	return nil
}

func (t *TransportOperator) AddDriver(driverID primitive.ObjectID) error {
	t.DriverIDs = append(t.DriverIDs, driverID)
	update := bson.M{"$set": bson.M{"driverIDs": t.DriverIDs}}
	_, err := db.Collection("transportOperator").UpdateOne(context.TODO(), bson.M{"_id": t.ID}, update)
	return err
}

func FindTransportOperatorsByDriverID(filter bson.M) ([]TransportOperator, error) {

	tos := []TransportOperator{}
	filter["isVerified"] = true

	cursor, err := db.Collection("transportOperator").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &tos); err != nil {
		return nil, err
	}
	return tos, nil
}
