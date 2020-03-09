package model

import (
	"context"
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	valid "github.com/asaskevich/govalidator"

	locModel "github.com/chadhao/logit/modules/location/model"
	"github.com/chadhao/logit/utils"
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
	Address locModel.Address `bson:"address,omitempty" json:"address,omitempty" valid:"-"`
	Coors   locModel.Coors   `bson:"coors,omitempty" json:"coors,omitempty" valid:"-"`
}

func (l *Location) equal(o *Location) bool {
	return l.Address == o.Address
}

// fillFull 若其中一项不完整，则用另外一项查找并补完
func (l *Location) fillFull() (err error) {
	if l.Address != "" && !l.Coors.EmptyCoors() {
		return
	}
	if l.Address == "" && l.Coors.EmptyCoors() {
		return errors.New("at least one field requried")
	}
	if l.Address == "" {
		l.Address, err = l.Coors.GetAddrFromCoors()
	} else {
		l.Coors, err = l.Address.GetCoorsFromAddr()
	}
	return
}

// Record 记录
type Record struct {
	ID            primitive.ObjectID `bson:"_id" json:"id" valid:"-"`
	DriverID      primitive.ObjectID `bson:"driverID" json:"driverID" valid:"required"`
	Type          Type               `bson:"type" json:"type" valid:"required"`
	Time          time.Time          `bson:"time" json:"time" valid:"required"`
	Duration      time.Duration      `bson:"duration" json:"duration" valid:"required"`
	StartLocation Location           `bson:"startLocation" json:"startLocation" valid:"required"`
	EndLocation   Location           `bson:"endLocation," json:"endLocation" valid:"required"`
	VehicleID     primitive.ObjectID `bson:"vehicleID" json:"vehicleID" valid:"required"`
	StartMileAge  *float64           `bson:"startDistance,omitempty" json:"startDistance,omitempty" valid:"-"`
	EndMileAge    *float64           `bson:"endDistance,omitempty" json:"endDistance,omitempty" valid:"-"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
	ClientTime    *time.Time         `bson:"clientTime,omitempty" json:"clientTime,omitempty" valid:"-"`
	DeletedAt     *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty" valid:"-"`
}

// Add 记录添加
func (r *Record) Add() (err error) {

	lastRec, err := getLastestRecord(r.DriverID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if err = r.beforeAdd(lastRec); err != nil {
		return
	}

	// 数据库添加记录
	if _, err = recordCollection.InsertOne(context.TODO(), r); err != nil {
		return
	}
	return
}

// Delete 记录删除
func (r *Record) Delete() error {
	switch {
	case r.DeletedAt != nil:
		return errors.New("record has already been deleted")
	case !r.isLatestRecord():
		return errors.New("record is not the lastest one")
	}
	// 数据库记录添加删除标记
	update := bson.M{"$set": bson.M{"deletedAt": time.Now()}}
	if _, err := recordCollection.UpdateOne(context.TODO(), bson.M{"_id": r.ID}, update); err != nil {
		return err
	}
	return nil
}

func (r *Record) beforeAdd(lastRec *Record) error {
	if (Record{}) != *lastRec {
		if lastRec.Type == r.Type {
			return errors.New("work type conflict with last record")
		}
		if lastRec.Time.After(r.Time) {
			return errors.New("time conflict with last record")
		}
		if !r.StartLocation.equal(&lastRec.EndLocation) {
			return errors.New("location not match")
		}
		if r.StartMileAge != nil && !utils.AlmostEqual(*r.StartMileAge, *lastRec.EndMileAge) {
			return errors.New("mileage not match")
		}
		if math.Abs(lastRec.Time.Add(r.Duration).Sub(r.Time).Seconds()) > 10 {
			return errors.New("time and duration not match")
		}
	}

	if err := r.EndLocation.fillFull(); err != nil {
		return err
	}

	return r.valid()
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
	lastest, err := getLastestRecord(r.DriverID)
	if err != nil {
		return false
	}
	return lastest.ID == r.ID
}

func getLastestRecord(driverID primitive.ObjectID) (*Record, error) {
	lastRec := new(Record)
	opts := options.FindOne().SetSort(bson.D{{Key: "time", Value: -1}})
	err := recordCollection.FindOne(context.TODO(), bson.M{"driverID": driverID, "deletedAt": nil}, opts).Decode(lastRec)
	return lastRec, err
}

// GetRecord 通过id获取记录
func GetRecord(id primitive.ObjectID) (*Record, error) {
	r := new(Record)
	err := recordCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(r)
	return r, err
}

// GetRecords 获取用户时间段内的记录
func GetRecords(driverID primitive.ObjectID, from, to time.Time, getDeleted bool) ([]Record, error) {
	records := []Record{}
	filter := bson.M{
		"driverID": driverID,
		"time":     bson.M{"$gte": from, "$lte": to},
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

// Records .
type Records []Record

// SyncAdd 批量上传添加
func (rs Records) SyncAdd() (err error) {

	lastRec, err := getLastestRecord(rs[0].DriverID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	l := len(rs)
	for i := 0; i < l; i++ {
		if err = rs[i].beforeAdd(lastRec); err != nil {
			return
		}
		lastRec = &rs[i]
	}

	// 数据库添加记录
	rsI := make([]interface{}, l)
	for i := range rs {
		rsI[i] = rs[i]
	}
	if _, err = recordCollection.InsertMany(context.TODO(), rsI); err != nil {
		return
	}
	return
}
