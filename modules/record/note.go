package record

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// INote 笔记接口
type INote interface{}

// Note 笔记
type Note struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Type      NoteType           `bson:"noteType" json:"noteType"`
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// SystemNote 系统笔记
type SystemNote struct {
	Note
}

// OtherWork 其它笔记
type OtherWork struct {
	Note
}

// ModificationNote 人为修改笔记
type ModificationNote struct {
	Note
	By primitive.ObjectID `bson:"by" json:"by"`
}

// TripNote 行程笔记
type TripNote struct {
	Note
	TransportOperatorID *primitive.ObjectID `bson:"transportOperatorID,omitempty" json:"transportOperatorID,omitempty"`
	StartTime           time.Time           `bson:"startTime" json:"startTime"`
	EndTime             *time.Time          `bson:"endTime,omitempty" json:"endTime,omitempty"`
	StartLocation       Location            `bson:"startLocation" json:"startLocation"`
	EndLocation         *Location           `bson:"endLocation,omitempty" json:"endLocation,omitempty"`
}
