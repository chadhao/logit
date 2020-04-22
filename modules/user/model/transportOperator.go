package model

import (
	"context"
	"errors"
	"time"

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

func (t *TransportOperator) Update() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.D{{"_id", t.ID}}
	tBson, err := bson.Marshal(t)
	if err != nil {
		return err
	}

	result, err := db.Collection("transportOperator").ReplaceOne(ctx, filter, tBson)
	if err != nil {
		return nil
	}
	if result.MatchedCount != 1 {
		return errors.New("transportOperator not updated")
	}
	return nil
}

func (f *TransportOperator) Filter(notVerifiedInclude bool) ([]TransportOperator, error) {

	tos := []TransportOperator{}

	filter := bson.M{}
	if !notVerifiedInclude {
		filter["isVerified"] = true
	}

	if len(f.LicenseNumber) > 0 {
		filter["licenseNumber"] = primitive.Regex{Pattern: f.LicenseNumber, Options: "i"}
	}
	if len(f.Name) > 0 {
		filter["name"] = primitive.Regex{Pattern: f.Name, Options: "i"}
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

func (t *TransportOperator) AddIdentity(userID primitive.ObjectID, identity TOIdentity, contact *string) (*TransportOperatorIdentity, error) {

	toI := &TransportOperatorIdentity{
		ID:                  primitive.NewObjectID(),
		TransportOperatorID: t.ID,
		CreatedAt:           time.Now(),
	}

	switch {
	case !t.IsVerified:
		tos, err := toI.Filter()
		if len(tos) != 0 {
			return nil, errors.New("transport operator need to be verified")
		}
		if err != nil {
			return nil, err
		}
	case t.IsCompany:
		if contact == nil && identity != TO_DRIVER {
			return nil, errors.New("contact is required")
		}
		toI.Contact = contact
	case !t.IsCompany:
		tos, err := toI.Filter()
		if len(tos) > 1 {
			return nil, errors.New("can only have one super")
		}
		if err != nil {
			return nil, err
		}
	}

	toI.UserID = userID
	toI.Identity = identity

	if err := toI.create(); err != nil {
		return nil, err
	}
	return toI, nil
}

func (t *TransportOperatorIdentity) Filter() ([]TransportOperatorIdentity, error) {
	tos := []TransportOperatorIdentity{}
	filter := bson.M{}
	if !t.UserID.IsZero() {
		filter["userID"] = t.UserID
	}
	if !t.TransportOperatorID.IsZero() {
		filter["transportOperatorID"] = t.TransportOperatorID
	}
	if t.Identity != "" {
		filter["identity"] = t.Identity
	}

	cursor, err := db.Collection("transportOperatorIdentity").Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &tos); err != nil {
		return nil, err
	}
	return tos, nil
}

func (t *TransportOperatorIdentity) Find() error {
	err := db.Collection("transportOperatorIdentity").FindOne(context.TODO(), bson.M{"_id": t.ID}).Decode(t)
	return err
}

func (t *TransportOperatorIdentity) UpdateContact(contact string) error {
	update := bson.M{"$set": bson.M{"contact": contact}}
	_, err := db.Collection("transportOperatorIdentity").UpdateOne(context.TODO(), bson.M{"_id": t.ID}, update)
	return err
}

func IsIdentity(uid, transportOperatorID primitive.ObjectID, identities []TOIdentity) bool {

	filter := bson.M{
		"userID":              uid,
		"transportOperatorID": transportOperatorID,
		"identity":            bson.M{"$in": identities},
	}

	if count, _ := db.Collection("transportOperatorIdentity").CountDocuments(context.TODO(), filter); count > 0 {
		return true
	}
	return false
}

func HasAccessTo(adminIDPlus, driverID, transportOperatorID primitive.ObjectID) bool {
	adminTos := []TransportOperatorIdentity{}
	adminFilter := bson.M{
		"userID":              adminIDPlus,
		"transportOperatorID": transportOperatorID,
		"identity":            bson.M{"$in": []TOIdentity{TO_ADMIN, TO_SUPER}},
	}
	cursor, err := db.Collection("transportOperatorIdentity").Find(context.TODO(), adminFilter)
	if err != nil {
		return false
	}
	if err = cursor.All(context.TODO(), &adminTos); err != nil {
		return false
	}

	driverTos := []TransportOperatorIdentity{}
	driverFilter := bson.M{
		"userID":              driverID,
		"transportOperatorID": transportOperatorID,
		"identity":            TO_DRIVER,
	}
	cursor, err = db.Collection("transportOperatorIdentity").Find(context.TODO(), driverFilter)
	if err != nil {
		return false
	}
	if err = cursor.All(context.TODO(), &driverTos); err != nil {
		return false
	}

	if len(adminTos) > 0 && len(driverTos) > 0 {
		return true
	}
	return false
}

func (t *TransportOperatorIdentity) create() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if t.exists() {
		return errors.New("This identity exists")
	}

	toBson, err := bson.Marshal(t)
	if err != nil {
		return err
	}

	if _, err := db.Collection("transportOperatorIdentity").InsertOne(ctx, toBson); err != nil {
		return err
	}

	return nil
}

func (t *TransportOperatorIdentity) Delete() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := db.Collection("transportOperatorIdentity").DeleteOne(ctx, bson.M{"_id": t.ID}); err != nil {
		return err
	}
	return nil
}

func (t *TransportOperatorIdentity) exists() bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conditions := primitive.A{
		bson.D{{"userID", t.UserID}},
		bson.D{{"transportOperatorID", t.TransportOperatorID}},
		bson.D{{"identity", t.Identity}},
		bson.D{{"deletedAt", nil}},
	}

	filter := bson.D{{"$and", conditions}}

	if count, _ := db.Collection("transportOperatorIdentity").CountDocuments(ctx, filter); count > 0 {
		return true
	}

	return false
}

// GetIdentitiesByUserIDs 批量通过userIDs获取他们相关的identity信息
func GetIdentitiesByUserIDs(uids []primitive.ObjectID) map[primitive.ObjectID][]TransportOperatorIdentity {
	m := make(map[primitive.ObjectID][]TransportOperatorIdentity)
	var tois []TransportOperatorIdentity
	filter := bson.M{
		"userID": bson.M{"$in": uids},
	}
	cursor, _ := db.Collection("transportOperatorIdentity").Find(context.TODO(), filter)
	cursor.All(context.TODO(), &tois)

	for _, toi := range tois {
		m[toi.UserID] = append(m[toi.UserID], toi)
	}
	return m
}
