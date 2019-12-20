package record

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	RecordType int
	NoteType   int
	HrTime     int
)

const (
	WORK RecordType = iota + 1
	REST
	SYSTEMNOTE NoteType = iota + 1
	MODIFICATIONNOTE
	TRIPNOTE
	OTHERWORKNOTE
	HR0D5 HrTime = iota
	HR5D5
	HR7
	HR10
	HR13
	HR24
	HR70
)

func (t HrTime) GetHrs() float64 {
	switch t {
	case HR0D5:
		return 0.5
	case HR5D5:
		return 5.5
	case HR7:
		return 7
	case HR10:
		return 10
	case HR13:
		return 13
	case HR24:
		return 24
	case HR70:
		return 70
	default:
		return 0
	}
}

// Location .
type Location struct {
	Address string     `bson:"address" json:"address"`
	Coors   [2]float64 `bson:"coors,omitempty" json:"coors,omitempty"`
}

// Record .
type Record struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	RecordType    RecordType         `bson:"recordType" json:"recordType"`
	StartTime     time.Time          `bson:"startTime" json:"startTime"`
	EndTime       time.Time          `bson:"endTime,omitempty" json:"endTime,omitempty"`
	StartLocation Location           `bson:"startLocation" json:"startLocation"`
	EndLocation   Location           `bson:"endLocation,omitempty" json:"endLocation,omitempty"`
	VehicleID     primitive.ObjectID `bson:"vehicleID" json:"vehicleID"`
	StartMileAge  *float64           `bson:"startDistance,omitempty" json:"startDistance,omitempty"`
	EndMileAge    *float64           `bson:"endDistance,omitempty" json:"endDistance,omitempty"`
	Notes         []Note             `bson:"notes,omitempty" json:"notes,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	DeletedAt     *time.Time         `bson:"isDeleted,omitempty" json:"isDeleted,omitempty"`
}

func (r *Record) AddNote()      {}
func (r *Record) Add()          {}
func (r *Record) Delete()       {}
func (r *Record) beforeAdd()    {}
func (r *Record) beforeDelete() {}
func (r *Record) valid()        {}
