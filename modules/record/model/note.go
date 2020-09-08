package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Note 笔记
type Note struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	RecordID  primitive.ObjectID `bson:"recordID" json:"recordID"`
	Type      NoteType           `bson:"noteType" json:"noteType"`
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// SystemNote 系统笔记
type SystemNote struct {
	Note `bson:",inline"`
}

// Add 系统笔记添加到数据库
func (sn *SystemNote) Add() error {
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), sn); err != nil {
		return err
	}
	return nil
}

// OtherWorkNote 其它笔记
type OtherWorkNote struct {
	Note `bson:",inline"`
}

// Add 其它笔记添加到数据库
func (own *OtherWorkNote) Add() error {
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), own); err != nil {
		return err
	}
	return nil
}

// ModificationNote 人为修改笔记
type ModificationNote struct {
	Note `bson:",inline"`
	By   primitive.ObjectID `bson:"by" json:"by"`
}

// Add 人为修改笔记添加到数据库
func (mn *ModificationNote) Add() error {
	// 数据库添加记录
	if _, err := noteCollection.InsertOne(context.TODO(), mn); err != nil {
		return err
	}
	return nil
}

// TripNote 行程笔记
type TripNote struct {
	Note                `bson:",inline"`
	TransportOperatorID primitive.ObjectID `bson:"transportOperatorID" json:"transportOperatorID"`
	StartTime           time.Time          `bson:"startTime" json:"startTime"`
	EndTime             time.Time          `bson:"endTime" json:"endTime"`
	StartLocation       Location           `bson:"startLocation" json:"startLocation"`
	EndLocation         Location           `bson:"endLocation" json:"endLocation"`
}

// Add 行程笔记添加到数据库
func (tn *TripNote) Add() error {
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
