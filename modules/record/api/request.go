package api

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/record/model"
)

// reqRecords 请求获取记录
type reqRecords struct {
	DriverID primitive.ObjectID `json:"driverID" query:"driverID" valid:"required"`
	From     time.Time          `json:"from" query:"from" valid:"required"`
	To       time.Time          `json:"to" query:"to" valid:"optional"`
}

func (reqR *reqRecords) valid() error {
	if _, err := valid.ValidateStruct(reqR); err != nil {
		return err
	}
	if reqR.To.IsZero() {
		reqR.To = time.Now()
	}
	if reqR.From.After(reqR.To) {
		return errors.New("times order is wrong")
	}
	return nil
}

// getRecords 获取指定时间范围内的记录
func (reqR *reqRecords) getRecords() ([]*respRecord, error) {
	if err := reqR.valid(); err != nil {
		return nil, err
	}
	// 获取记录
	records, err := model.GetRecords(reqR.DriverID, reqR.From, reqR.To, false)
	if err != nil {
		return nil, err
	}
	// 获取记录下的笔记
	recordIDs := []primitive.ObjectID{}
	for _, v := range records {
		recordIDs = append(recordIDs, v.ID)
	}
	notesMap, err := model.GetNotesByRecordIDs(recordIDs)
	if err != nil {
		return nil, err
	}
	// 拼装返回
	respRecords := []*respRecord{}
	for _, v := range records {
		respRecords = append(respRecords, &respRecord{
			Record: v,
			Notes:  notesMap[v.ID],
		})
	}

	return respRecords, nil
}

// reqRecord 请求获取记录
type reqRecord struct {
	ID primitive.ObjectID `json:"id" query:"id" valid:"required"`
}

// getRecord 获取记录
func (reqRecord *reqRecord) getRecord() (*respRecord, error) {
	// 获取记录
	r, err := model.GetRecord(reqRecord.ID)
	if err != nil {
		return nil, err
	}
	// 获取记录下的笔记
	notesMap, err := model.GetNotesByRecordIDs([]primitive.ObjectID{reqRecord.ID})
	if err != nil {
		return nil, err
	}
	respRecord := &respRecord{
		Record: *r,
		Notes:  notesMap[reqRecord.ID],
	}
	// 拼装返回
	return respRecord, nil
}

// deleteRecord 删除记录
func (reqRecord *reqRecord) deleteRecord(userID primitive.ObjectID) error {
	// 获取记录
	r, err := model.GetRecord(reqRecord.ID)
	if err != nil {
		return err
	}
	if r.UserID != userID {
		return errors.New("no authorization")
	}

	return r.Delete()
}

// reqAddRecord 添加记录请求结构
type reqAddRecord struct {
	Type          model.Type      `json:"type" valid:"required"`
	StartTime     *time.Time      `json:"startTime,omitempty" valid:"optional"`
	EndTime       *time.Time      `json:"endTime,omitempty" valid:"optional"`
	StartLocation model.Location  `json:"startLocation" valid:"required"`
	EndLocation   *model.Location `json:"endLocation,omitempty" valid:"-"`
	StartMileAge  *float64        `json:"startDistance,omitempty" valid:"-"`
	EndMileAge    *float64        `json:"endDistance,omitempty" valid:"-"`
	ClientTime    *time.Time      `json:"clientTime,omitempty" valid:"optional"`
}

// Valid 添加记录请求结构验证
func (reqAddR *reqAddRecord) valid() error {
	if reqAddR.Type != model.WORK && reqAddR.Type != model.REST {
		return errors.New("no match type")
	}
	if _, err := valid.ValidateStruct(reqAddR); err != nil {
		return err
	}
	// 1. 时间检验
	if reqAddR.StartTime != nil && reqAddR.EndTime != nil {
		if reqAddR.StartTime.After(*reqAddR.EndTime) {
			return errors.New("startTime should be before endTime")
		}
	}
	if reqAddR.StartTime != nil && reqAddR.StartTime.After(time.Now()) {
		return errors.New("cannot add future time to startTime")
	}
	if reqAddR.EndTime != nil && reqAddR.EndTime.After(time.Now()) {
		return errors.New("cannot add future time to endTime")
	}

	// 2. 若公里数不为空时的检验
	if reqAddR.StartMileAge != nil && reqAddR.EndMileAge != nil && *reqAddR.StartMileAge > *reqAddR.EndMileAge {
		return errors.New("startMileAge should be less than endMileAge")
	}
	return nil
}

// constructToRecord 将reqAddRecord构造为Record
func (reqAddR *reqAddRecord) constructToRecord(userID, vehicleID primitive.ObjectID) (*model.Record, error) {
	if err := reqAddR.valid(); err != nil {
		return nil, err
	}
	t := time.Now()
	if reqAddR.StartTime != nil {
		t = *reqAddR.StartTime
	}
	r := &model.Record{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		Type:          reqAddR.Type,
		StartTime:     t,
		EndTime:       reqAddR.EndTime,
		StartLocation: reqAddR.StartLocation,
		EndLocation:   reqAddR.EndLocation,
		VehicleID:     vehicleID,
		StartMileAge:  reqAddR.StartMileAge,
		EndMileAge:    reqAddR.EndMileAge,
		ClientTime:    reqAddR.ClientTime,
		CreatedAt:     t,
	}
	return r, nil
}

type reqAddNote struct {
	NoteType model.NoteType     `json:"noteType" valid:"required"`
	RecordID primitive.ObjectID `json:"recordID" valid:"required"`
	Comment  string             `json:"comment" valid:"optional"`
}

// valid 添加笔记验证
func (r *reqAddNote) valid() error {
	_, err := valid.ValidateStruct(r)
	return err
}

// constructToSystemNote reqAddNote
func (r *reqAddNote) constructToSystemNote() (*model.SystemNote, error) {
	if err := r.valid(); err != nil {
		return nil, err
	}
	sn := &model.SystemNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  r.RecordID,
			Type:      r.NoteType,
			Comment:   r.Comment,
			CreatedAt: time.Now(),
		},
	}
	return sn, nil
}

// constructToOtherWorkNote 将reqAddNote构造为OtherWorkNote
func (r *reqAddNote) constructToOtherWorkNote() (*model.OtherWorkNote, error) {
	if err := r.valid(); err != nil {
		return nil, err
	}
	own := &model.OtherWorkNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  r.RecordID,
			Type:      r.NoteType,
			Comment:   r.Comment,
			CreatedAt: time.Now(),
		},
	}
	return own, nil
}

// constructToModificationNote 将reqAddNote构造为ModificationNote
func (r *reqAddNote) constructToModificationNote(by primitive.ObjectID) (*model.ModificationNote, error) {
	if err := r.valid(); err != nil {
		return nil, err
	}
	mn := &model.ModificationNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  r.RecordID,
			Type:      r.NoteType,
			Comment:   r.Comment,
			CreatedAt: time.Now(),
		},
		By: by,
	}
	return mn, nil
}

func (r *reqAddNote) isUsersRecord(userID primitive.ObjectID) bool {
	rec, err := model.GetRecord(r.RecordID)
	if err != nil {
		return false
	}
	return rec.UserID == userID
}
