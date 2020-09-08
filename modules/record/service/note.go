package service

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/record/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// INote 笔记接口
type INote interface {
	Add() error
}

type (
	// AddNoteInput 添加笔记参数
	AddNoteInput struct {
		RecordID primitive.ObjectID `json:"recordID" valid:"required"`
		Type     model.NoteType     `json:"noteType" valid:"required"`
		Comment  string             `json:"comment" valid:"-"`

		// 以下为tripNote所需参数
		TransportOperatorID primitive.ObjectID `json:"transportOperatorID" valid:"-"`
		StartTime           time.Time          `json:"startTime" valid:"-"`
		EndTime             time.Time          `json:"endTime" valid:"-"`
		StartLocation       model.Location     `json:"startLocation" valid:"-"`
		EndLocation         model.Location     `json:"endLocation" valid:"-"`

		// 以下为ModificationNote所需参数
		By primitive.ObjectID `json:"by" valid:"-"`
	}
	// AddNoteOutput 添加笔记返回参数
	AddNoteOutput struct {
		INote
	}
)

// IsDriversRecord 检查该笔记是否是传入司机的记录
func (n *AddNoteInput) IsDriversRecord(driverID primitive.ObjectID) bool {
	rec, err := model.GetRecord(n.RecordID)
	if err != nil {
		return false
	}
	return rec.DriverID == driverID
}

// toSystemNote 将AddNoteInput构造为SystemNote
func (n *AddNoteInput) toSystemNote() (*model.SystemNote, error) {
	sn := &model.SystemNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  n.RecordID,
			Type:      n.Type,
			Comment:   n.Comment,
			CreatedAt: time.Now(),
		},
	}
	return sn, nil
}

// toOtherWorkNote 将AddNoteInput构造为OtherWorkNote
func (n *AddNoteInput) toOtherWorkNote() (*model.OtherWorkNote, error) {
	own := &model.OtherWorkNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  n.RecordID,
			Type:      n.Type,
			Comment:   n.Comment,
			CreatedAt: time.Now(),
		},
	}
	return own, nil
}

// toModificationNote 将reqAddNote构造为ModificationNote
func (n *AddNoteInput) toModificationNote(by primitive.ObjectID) (*model.ModificationNote, error) {
	mn := &model.ModificationNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  n.RecordID,
			Type:      n.Type,
			Comment:   n.Comment,
			CreatedAt: time.Now(),
		},
		By: by,
	}
	return mn, nil
}

// ToTripNote 将reqAddNote构造为TripNote
func (n *AddNoteInput) toTripNote() (*model.TripNote, error) {
	if n.TransportOperatorID.IsZero() {
		return nil, errors.New("transportOperatorID is required")
	}
	if n.StartTime.IsZero() {
		return nil, errors.New("startTime is required")
	}
	if n.EndTime.IsZero() {
		return nil, errors.New("endTime is required")
	}
	if (model.Location{}) == n.StartLocation {
		return nil, errors.New("startLocation is required")
	}
	if (model.Location{}) == n.EndLocation {
		return nil, errors.New("endLocation is required")
	}
	tn := &model.TripNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  n.RecordID,
			Type:      n.Type,
			Comment:   n.Comment,
			CreatedAt: time.Now(),
		},
		TransportOperatorID: n.TransportOperatorID,
		StartTime:           n.StartTime,
		EndTime:             n.EndTime,
		StartLocation:       n.StartLocation,
		EndLocation:         n.EndLocation,
	}
	return tn, nil
}

// AddNote 添加笔记(除了SystemNote)
func AddNote(in *AddNoteInput) (*AddNoteOutput, error) {
	var (
		note INote
		err  error
	)

	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}

	switch in.Type {
	case model.OTHERWORKNOTE:
		note, err = in.toOtherWorkNote()
		if err != nil {
			return nil, err
		}
	case model.MODIFICATIONNOTE:
		note, err = in.toModificationNote(in.By)
		if err != nil {
			return nil, err
		}
	case model.TRIPNOTE:
		note, err = in.toTripNote()
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("no match noteType")
	}

	if err = note.Add(); err != nil {
		return nil, err
	}
	return &AddNoteOutput{note}, nil
}
