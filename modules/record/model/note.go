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
}

// Note 笔记
type Note struct {
	ID        primitive.ObjectID `bson:"_id" json:"id" valid:"-"`
	RecordID  primitive.ObjectID `bson:"recordID" json:"recordID" valid:"required"`
	Type      NoteType           `bson:"noteType" json:"noteType" valid:"required"`
	Comment   string             `bson:"comment" json:"comment" valid:"-"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
}

func (n *Note) valid() error {
	_, err := valid.ValidateStruct(n)
	return err
}

// SystemNote 系统笔记
type SystemNote struct {
	Note `bson:",inline" valid:"-"`
}

// Add 系统笔记添加到数据库
func (sn *SystemNote) Add() error {
	if _, err := valid.ValidateStruct(sn); err != nil {
		return err
	}
	if err := sn.Note.valid(); err != nil {
		return err
	}
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), sn); err != nil {
		return err
	}
	return nil
}

// OtherWorkNote 其它笔记
type OtherWorkNote struct {
	Note `bson:",inline" valid:"-"`
}

// Add 其它笔记添加到数据库
func (own *OtherWorkNote) Add() error {
	if _, err := valid.ValidateStruct(own); err != nil {
		return err
	}
	if err := own.Note.valid(); err != nil {
		return err
	}
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), own); err != nil {
		return err
	}
	return nil
}

// ModificationNote 人为修改笔记
type ModificationNote struct {
	Note `bson:",inline" valid:"-"`
	By   primitive.ObjectID `bson:"by" json:"by" valid:"required"`
}

// Add 人为修改笔记添加到数据库
func (mn *ModificationNote) Add() error {
	if _, err := valid.ValidateStruct(mn); err != nil {
		return err
	}
	if err := mn.Note.valid(); err != nil {
		return err
	}
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), mn); err != nil {
		return err
	}
	return nil
}

// TripNote 行程笔记
type TripNote struct {
	Note                `bson:",inline" valid:"-"`
	TransportOperatorID primitive.ObjectID `bson:"transportOperatorID" json:"transportOperatorID" valid:"required"`
	StartTime           time.Time          `bson:"startTime" json:"startTime" valid:"required"`
	EndTime             time.Time          `bson:"endTime" json:"endTime" valid:"required"`
	StartLocation       Location           `bson:"startLocation" json:"startLocation" valid:"required"`
	EndLocation         Location           `bson:"endLocation" json:"endLocation" valid:"required"`
}

// Add 行程笔记添加到数据库
func (tn *TripNote) Add() error {
	if _, err := valid.ValidateStruct(tn); err != nil {
		return err
	}
	if err := tn.Note.valid(); err != nil {
		return err
	}

	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), tn); err != nil {
		return err
	}
	return nil
}

type (
	// DifNote 不同的笔记
	DifNote bson.M
	// DifNotes 不同的笔记集合
	DifNotes []DifNote
)

func (d DifNote) getRecordID() primitive.ObjectID {
	return d["recordID"].(primitive.ObjectID)
}

// GetNotesByRecordIDs 获取recordIDs相对应的记录并以key为recordID,value为所对应Notes的map结构返回
func GetNotesByRecordIDs(recordIDs []primitive.ObjectID) (map[primitive.ObjectID]DifNotes, error) {

	notes := DifNotes{}

	cursor, err := noteCollection.Find(context.TODO(), bson.M{"recordID": bson.M{"$in": recordIDs}})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &notes); err != nil {
		return nil, err
	}

	notesMap := make(map[primitive.ObjectID]DifNotes)
	for _, v := range notes {
		notesMap[v.getRecordID()] = append(notesMap[v.getRecordID()], v)
	}

	return notesMap, nil
}
