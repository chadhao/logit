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

type TransportOperatorFilter struct {
	DriverID primitive.ObjectID `json:"driverID"`
	SuperID  primitive.ObjectID `json:"superID"`
	AdminID  primitive.ObjectID `json:"adminID"`
	Name     string             `json:"name"`
}

func (f *TransportOperatorFilter) Find() ([]TransportOperator, error) {

	tos := []TransportOperator{}

	filter := bson.M{
		// "isVerified": true,
	}

	switch {
	case !f.DriverID.IsZero():
		filter["driverIDs"] = f.DriverID
	case !f.SuperID.IsZero():
		filter["superIDs"] = f.SuperID
	case !f.AdminID.IsZero():
		filter["adminIDs"] = f.AdminID
	case len(f.Name) > 0:
		filter["name"] = bson.M{"$regex": "(?i)" + f.Name}
	default:
		return nil, errors.New("filter field is required")
	}

	cursor, err := db.Collection("transportOperator").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &tos); err != nil {
		return nil, err
	}
	return tos, nil
}

func (f *TransportOperatorFilter) FindTransportOperatorsRelatedToUser() (drivers []TransportOperator, supers []TransportOperator, admins []TransportOperator, err error) {

	drivers = []TransportOperator{}
	supers = []TransportOperator{}
	admins = []TransportOperator{}

	filter := bson.M{
		"driverIDs": f.DriverID,
	}
	cursor, err := db.Collection("transportOperator").Find(context.TODO(), filter)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = cursor.All(context.TODO(), &drivers); err != nil {
		return nil, nil, nil, err
	}

	filter = bson.M{
		"superIDs": f.SuperID,
	}
	cursor, err = db.Collection("transportOperator").Find(context.TODO(), filter)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = cursor.All(context.TODO(), &supers); err != nil {
		return nil, nil, nil, err
	}

	filter = bson.M{
		"adminIDs": f.AdminID,
	}
	cursor, err = db.Collection("transportOperator").Find(context.TODO(), filter)
	if err != nil {
		return nil, nil, nil, err
	}
	if err = cursor.All(context.TODO(), &admins); err != nil {
		return nil, nil, nil, err
	}
	return drivers, supers, admins, nil
}
