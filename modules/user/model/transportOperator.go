package model

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (t *TransportOperator) Delete() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	filter := bson.D{{"_id", t.ID}}

	if _, err := db.Collection("transportOperator").DeleteOne(ctx, filter); err != nil {
		return nil
	}
	return nil

}

func (f *TransportOperator) Filter(driverOrigin bool) ([]TransportOperator, error) {

	tos := []TransportOperator{}

	filter := bson.M{}
	if driverOrigin {
		filter["isVerified"] = true
		filter["isCompany"] = true
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

	// to super, to admin 用户在公司情况下需要填入contact信息
	if t.IsCompany {
		if contact == nil && identity != TO_DRIVER {
			return nil, errors.New("contact is required")
		}
		toI.Contact = contact
	} else { // 自雇形式
		// 只能添加自己
		// 添加为super, 只能存在一个super
		if identity == TO_SUPER {
			toI.Identity = TO_SUPER
			if tos, _ := toI.Filter(); len(tos) >= 1 {
				return nil, errors.New("can only have one super")
			}
		} else if identity == TO_DRIVER { // 只能将自己添加为driver
			toI.Identity = TO_SUPER
			tos, _ := toI.Filter()
			if len(tos) == 0 {
				return nil, errors.New("no super found")
			}
			if tos[0].UserID != userID {
				return nil, errors.New("to super id and this id not match")
			}

			toI.Identity = TO_DRIVER
			if tos, _ = toI.Filter(); len(tos) != 0 {
				return nil, errors.New("has one driver already")
			}
		} else {
			return nil, errors.New("identity can only be to super or driver")
		}
	}

	toI.UserID = userID
	toI.Identity = identity

	if err := toI.create(); err != nil {
		return nil, err
	}
	return toI, nil
}

func (t *TransportOperatorIdentity) Filter() ([]TransportOperatorIdentityDetail, error) {
	tos := []TransportOperatorIdentityDetail{}
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
	query := mongo.Pipeline{
		bson.D{
			{"$lookup", bson.D{
				{"from", "transportOperator"},
				{"localField", "transportOperatorID"},
				{"foreignField", "_id"},
				{"as", "transportOperator"},
			}},
		},
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$unwind", "$transportOperator"},
		},
	}
	cursor, err := db.Collection("transportOperatorIdentity").Aggregate(context.TODO(), query)
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
func GetIdentitiesByUserIDs(uids []primitive.ObjectID) (map[primitive.ObjectID][]TransportOperatorIdentityDetail, error) {
	m := make(map[primitive.ObjectID][]TransportOperatorIdentityDetail)
	var tois []TransportOperatorIdentityDetail
	filter := bson.M{
		"userID": bson.M{"$in": uids},
	}

	query := mongo.Pipeline{
		bson.D{
			{"$lookup", bson.D{
				{"from", "transportOperator"},
				{"localField", "transportOperatorID"},
				{"foreignField", "_id"},
				{"as", "transportOperator"},
			}},
		},
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$unwind", "$transportOperator"},
		},
	}
	cursor, err := db.Collection("transportOperatorIdentity").Aggregate(context.TODO(), query)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &tois); err != nil {
		return nil, err
	}
	for _, toi := range tois {
		m[toi.UserID] = append(m[toi.UserID], toi)
	}
	return m, nil
}
