package record

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"

	valid "github.com/asaskevich/govalidator"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// Type 记录类型
	Type string
	// NoteType 笔记类型
	NoteType string
	// HrTime 重要时间点
	HrTime float64
)

const (
	// WORK 工作记录类型
	WORK Type = "work"
	// REST 休息记录类型
	REST Type = "rest"
)
const (
	// SYSTEMNOTE 系统笔记类型
	SYSTEMNOTE NoteType = "system"
	// MODIFICATIONNOTE 人为修改笔记类型
	MODIFICATIONNOTE NoteType = "modification"
	// TRIPNOTE 行程笔记类型
	TRIPNOTE NoteType = "trip"
	// OTHERWORKNOTE 其它笔记类型
	OTHERWORKNOTE NoteType = "others"
)

const (
	// HR0D5 0.5小时
	HR0D5 HrTime = 0.5
	// HR5D5 5.5小时
	HR5D5 HrTime = 5.5
	// HR7 7.5小时
	HR7 HrTime = 7
	// HR10 10小时
	HR10 HrTime = 10
	// HR13 13小时
	HR13 HrTime = 13
	// HR24 24小时
	HR24 HrTime = 24
	// HR70 70小时
	HR70 HrTime = 70
)

func (t HrTime) getHrs() float64 {
	return float64(t)
}

// Location 位置信息
type Location struct {
	Address string     `bson:"address" json:"address"`
	Coors   [2]float64 `bson:"coors,omitempty" json:"coors,omitempty"`
}

// Record 记录
type Record struct {
	ID            primitive.ObjectID `bson:"_id" json:"id" valid:"required"`
	UserID        primitive.ObjectID `bson:"userID" json:"userID" valid:"required"`
	Type          Type               `bson:"type" json:"type" valid:"required"`
	StartTime     time.Time          `bson:"startTime" json:"startTime" valid:"required"`
	EndTime       time.Time          `bson:"endTime,omitempty" json:"endTime,omitempty" valid:"-"`
	StartLocation Location           `bson:"startLocation" json:"startLocation" valid:"required"`
	EndLocation   Location           `bson:"endLocation,omitempty" json:"endLocation,omitempty" valid:"-"`
	VehicleID     primitive.ObjectID `bson:"vehicleID" json:"vehicleID" valid:"required"`
	StartMileAge  *float64           `bson:"startDistance,omitempty" json:"startDistance,omitempty" valid:"-"`
	EndMileAge    *float64           `bson:"endDistance,omitempty" json:"endDistance,omitempty" valid:"-"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
	ClientTime    *time.Time         `bson:"clientTime,omitempty" json:"clientTime,omitempty" valid:"-"`
	DeletedAt     *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty" valid:"-"`
}

// Add 记录添加
func (r *Record) Add() (err error) {
	if err = r.beforeAdd(); err != nil {
		return
	}
	// 数据库添加记录
	if _, err = recordCollection.InsertOne(context.TODO(), r); err != nil {
		return
	}
	go r.afterAdd()
	return
}

// Delete 记录删除
func (r *Record) Delete(userID primitive.ObjectID) error {
	switch {
	case r.UserID != userID:
		return errors.New("no authorization")
	case r.DeletedAt != nil:
		return errors.New("record has already been deleted")
	case !r.isLatestRecord():
		return errors.New("record is not the lastest one")
	}
	// 数据库记录添加删除标记
	if _, err := recordCollection.UpdateOne(context.TODO(), bson.M{"_id": r.ID}, bson.M{"deletedAt": time.Now()}); err != nil {
		return err
	}
	return nil
}

func (r *Record) beforeAdd() error {
	return r.valid()
}

func (r *Record) afterAdd() error {
	// 1. 获取上一条记录，将上一条记录信息补充完整; endTime, endLocation不为空时候添加
	// 2. 若此条和上一条的startMileAge都不为空，上一条endMileAge不为空，则为上一条添加endMileAge
	lastRec, err := getLastestRecord(r.UserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return err
	}
	update := bson.M{
		"endTime":     r.StartTime,
		"endLocation": r.StartLocation,
	}
	if r.StartMileAge != nil && lastRec.StartMileAge != nil {
		update["endMileAge"] = r.StartMileAge
	}
	if _, err := recordCollection.UpdateOne(context.Background(), bson.M{"_id": lastRec.ID}, bson.M{"$set": update}); err != nil {
		return err
	}
	return nil
}

func (r *Record) valid() error {
	// 1. 验证记录结构是否完整
	if _, err := valid.ValidateStruct(r); err != nil {
		return err
	}
	// 2. 验证记录内容, TBC...
	return nil
}

func (r *Record) isLatestRecord() bool {
	lastest, err := getLastestRecord(r.UserID)
	if err != nil {
		return false
	}
	return lastest.ID == r.ID
}

func getLastestRecord(userID primitive.ObjectID) (*Record, error) {
	lastRec := new(Record)
	err := recordCollection.FindOne(context.TODO(), bson.M{"userID": userID, "$natural": -1, "deletedAt": nil}).Decode(lastRec)
	return lastRec, err
}

func getRecord(id primitive.ObjectID) (*Record, error) {
	r := new(Record)
	err := recordCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(r)
	return r, err
}

func getRecords(userID primitive.ObjectID, from, to time.Time, getDeleted bool) ([]Record, error) {
	records := []Record{}
	filter := bson.M{
		"userID": userID,
		"$gte":   bson.M{"startTime": from},
		"$lte":   bson.M{"endTime": to},
	}
	if !getDeleted {
		filter["deletedAt"] = nil
	}
	cursor, err := recordCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &records); err != nil {
		return nil, err
	}
	return records, nil
}
