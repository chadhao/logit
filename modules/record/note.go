package record

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type INote interface{}

// Note .
type Note struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Type      NoteType           `bson:"noteType" json:"noteType"`
	Comment   string             `bson:"comment" json:"comment"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

// SystemNote .
type SystemNote struct {
	Note
}

// OtherWork .
type OtherWork struct {
	Note
}

// ModificationNote .
type ModificationNote struct {
	Note
	By primitive.ObjectID `bson:"by" json:"by"`
}

// TripNote .
type TripNote struct {
	Note
	TransportOperatorID *primitive.ObjectID `bson:"transportOperatorID,omitempty" json:"transportOperatorID,omitempty"`
	StartTime           time.Time           `bson:"startTime" json:"startTime"`
	EndTime             *time.Time          `bson:"endTime,omitempty" json:"endTime,omitempty"`
	StartLocation       Location            `bson:"startLocation" json:"startLocation"`
	EndLocation         *Location           `bson:"endLocation,omitempty" json:"endLocation,omitempty"`
}

// SetTransportOperator .
func (t *TripNote) SetTransportOperator(id primitive.ObjectID) {
	t.TransportOperatorID = &id
}
