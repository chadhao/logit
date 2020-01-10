package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

func (d *Driver) Find() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filter bson.D
	if !d.Id.IsZero() {
		filter = bson.D{{"_id", d.Id}}
	} else if !d.UserId.IsZero() {
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
