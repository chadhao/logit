package model

import (
	"context"
	"time"

	valid "github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// INote 笔记接口
type INote interface {
	Add() error
	GetRecordID() primitive.ObjectID
}

// Note 笔记
type Note struct {
	ID        primitive.ObjectID `bson:"_id" json:"id" valid:"required"`
	RecordID  primitive.ObjectID `bson:"recordID" json:"recordID" valid:"required"`
	Type      NoteType           `bson:"noteType" json:"noteType" valid:"required"`
	Comment   string             `bson:"comment" json:"comment" valid:"-"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
}

// SystemNote 系统笔记
type SystemNote struct {
	Note `bson:",inline"`
}

// Add 系统笔记添加到数据库
func (sn *SystemNote) Add() error {
	if _, err := valid.ValidateStruct(sn); err != nil {
		return err
	}
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), sn); err != nil {
		return err
	}
	return nil
}

// GetRecordID 系统笔记获取RecordID
func (sn *SystemNote) GetRecordID() primitive.ObjectID {
	return sn.RecordID
}

// OtherWorkNote 其它笔记
type OtherWorkNote struct {
	Note `bson:",inline"`
}

// Add 其它笔记添加到数据库
func (own *OtherWorkNote) Add() error {
	if _, err := valid.ValidateStruct(own); err != nil {
		return err
	}
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), own); err != nil {
		return err
	}
	return nil
}

// GetRecordID 其它笔记获取RecordID
func (own *OtherWorkNote) GetRecordID() primitive.ObjectID {
	return own.RecordID
}

// ModificationNote 人为修改笔记
type ModificationNote struct {
	Note `bson:",inline"`
	By   primitive.ObjectID `bson:"by" json:"by" valid:"required"`
}

// Add 人为修改笔记添加到数据库
func (mn *ModificationNote) Add() error {
	if _, err := valid.ValidateStruct(mn); err != nil {
		return err
	}
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), mn); err != nil {
		return err
	}
	return nil
}

// GetRecordID 人为修改笔记获取RecordID
func (mn *ModificationNote) GetRecordID() primitive.ObjectID {
	return mn.RecordID
}

// TripNote 行程笔记
type TripNote struct {
	Note                `bson:",inline"`
	TransportOperatorID *primitive.ObjectID `bson:"transportOperatorID,omitempty" json:"transportOperatorID,omitempty" valid:"-"`
	StartTime           time.Time           `bson:"startTime" json:"startTime" valid:"required"`
	EndTime             time.Time           `bson:"endTime" json:"endTime" valid:"required"`
	StartLocation       Location            `bson:"startLocation" json:"startLocation" valid:"required"`
	EndLocation         Location            `bson:"endLocation" json:"endLocation" valid:"required"`
}

// Add 行程笔记添加到数据库
func (tn *TripNote) Add() error {
	if _, err := valid.ValidateStruct(tn); err != nil {
		return err
	}
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), tn); err != nil {
		return err
	}
	return nil
}

// GetRecordID 行程笔记获取RecordID
func (tn *TripNote) GetRecordID() primitive.ObjectID {
	return tn.RecordID
}

// GetNotes 获取recordIDs相对应的记录
func GetNotes(recordIDs []primitive.ObjectID) ([]INote, error) {
	notes := []INote{}
	cursor, err := noteCollection.Find(context.TODO(), bson.M{"$in": bson.M{"recordID": recordIDs}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &notes); err != nil {
		return nil, err
	}
	return notes, nil
}
