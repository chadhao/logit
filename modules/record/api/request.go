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
	From time.Time `json:"from" valid:"required"`
	To   time.Time `json:"to" valid:"-"`
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
func (reqR *reqRecords) getRecords(userID primitive.ObjectID) ([]*respRecord, error) {
	if err := reqR.valid(); err != nil {
		return nil, err
	}
	// 获取记录
	records, err := model.GetRecords(userID, reqR.From, reqR.To, false)
	if err != nil {
		return nil, err
	}
	// 获取记录下的笔记
	recordIDs := []primitive.ObjectID{}
	for _, v := range records {
		recordIDs = append(recordIDs, v.ID)
	}
	notes, err := model.GetNotesByRecordIDs(recordIDs)
	if err != nil {
		return nil, err
	}
	// 拼装返回
	respRecords := []*respRecord{}
	notesMap := make(map[primitive.ObjectID][]model.INote)
	for _, v := range notes {
		notesMap[v.GetRecordID()] = append(notesMap[v.GetRecordID()], v)
	}
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
	ID primitive.ObjectID `json:"id" valid:"required"`
}

// getRecord 获取记录
func (reqRecord *reqRecord) getRecord() (*respRecord, error) {
	// 获取记录
	r, err := model.GetRecord(reqRecord.ID)
	if err != nil {
		return nil, err
	}

	// 获取记录下的笔记
	notes, err := model.GetNotesByRecordIDs([]primitive.ObjectID{reqRecord.ID})
	if err != nil {
		return nil, err
	}
	respRecord := &respRecord{
		Record: *r,
		Notes:  notes,
	}
	// 拼装返回
	return respRecord, nil
}

// reqAddRecord 添加记录请求结构
type reqAddRecord struct {
	Type          model.Type     `json:"type" valid:"required"`
	StartTime     time.Time      `json:"startTime,omitempty" valid:"-"`
	EndTime       time.Time      `json:"endTime,omitempty" valid:"-"`
	StartLocation model.Location `json:"startLocation" valid:"required"`
	EndLocation   model.Location `json:"endLocation,omitempty" valid:"-"`
	StartMileAge  *float64       `json:"startDistance,omitempty" valid:"-"`
	EndMileAge    *float64       `json:"endDistance,omitempty" valid:"-"`
	ClientTime    *time.Time     `bson:"clientTime,omitempty" json:"clientTime,omitempty" valid:"-"`
}

// Valid 添加记录请求结构验证
func (reqAddR *reqAddRecord) valid() error {
	if _, err := valid.ValidateStruct(reqAddR); err != nil {
		return err
	}
	// 1. 时间检验
	if !reqAddR.StartTime.IsZero() && !reqAddR.EndTime.IsZero() {
		if reqAddR.StartTime.After(reqAddR.EndTime) {
			return errors.New("startTime should be before endTime")
		}
	}
	if reqAddR.EndTime.After(time.Now()) || reqAddR.StartTime.After(time.Now()) {
		return errors.New("cannot add future time")
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
	now := time.Now()
	r := &model.Record{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		Type:          reqAddR.Type,
		StartTime:     reqAddR.StartTime,
		EndTime:       reqAddR.EndTime,
		StartLocation: reqAddR.StartLocation,
		EndLocation:   reqAddR.EndLocation,
		VehicleID:     vehicleID,
		StartMileAge:  reqAddR.StartMileAge,
		EndMileAge:    reqAddR.EndMileAge,
		ClientTime:    reqAddR.ClientTime,
		CreatedAt:     now,
	}
	if r.StartTime.IsZero() {
		r.StartTime = now
	}
	return r, nil
}

// reqAddSystemNote 添加系统笔记
type reqAddSystemNote struct {
	RecordID primitive.ObjectID `json:"recordID" valid:"required"`
	Comment  string             `json:"comment" valid:"required"`
	Type     model.NoteType     `json:"noteType" valid:"required"`
}

// Valid 添加系统笔记验证
func (r *reqAddSystemNote) valid() error {
	_, err := valid.ValidateStruct(r)
	return err
}

// constructToSystemNote 将reqAddSystemNote构造为SystemNote
func (r *reqAddSystemNote) constructToSystemNote() (*model.SystemNote, error) {
	if err := r.valid(); err != nil {
		return nil, err
	}
	sn := &model.SystemNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  r.RecordID,
			Type:      r.Type,
			Comment:   r.Comment,
			CreatedAt: time.Now(),
		},
	}
	return sn, nil
}

// reqAddOtherWorkNote 添加其它笔记
type reqAddOtherWorkNote struct {
	RecordID primitive.ObjectID `json:"recordID" valid:"required"`
	Comment  string             `json:"comment" valid:"required"`
	Type     model.NoteType     `json:"noteType" valid:"required"`
}

// Valid 添加其它笔记验证
func (r *reqAddOtherWorkNote) valid() error {
	_, err := valid.ValidateStruct(r)
	return err
}

// constructToOtherWorkNote 将reqAddOtherWorkNote构造为OtherWorkNote
func (r *reqAddOtherWorkNote) constructToOtherWorkNote() (*model.OtherWorkNote, error) {
	if err := r.valid(); err != nil {
		return nil, err
	}
	own := &model.OtherWorkNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  r.RecordID,
			Type:      r.Type,
			Comment:   r.Comment,
			CreatedAt: time.Now(),
		},
	}
	return own, nil
}

// reqAddModificationNote 添加人为修改笔记
type reqAddModificationNote struct {
	RecordID primitive.ObjectID `json:"recordID" valid:"required"`
	Comment  string             `json:"comment" valid:"required"`
	Type     model.NoteType     `json:"noteType" valid:"required"`
}

// Valid 添加人为修改笔记
func (r *reqAddModificationNote) valid() error {
	_, err := valid.ValidateStruct(r)
	return err
}

// constructToModificationNote 将reqAddModificationNote构造为ModificationNote
func (r *reqAddModificationNote) constructToModificationNote(by primitive.ObjectID) (*model.ModificationNote, error) {
	if err := r.valid(); err != nil {
		return nil, err
	}
	mn := &model.ModificationNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  r.RecordID,
			Type:      r.Type,
			Comment:   r.Comment,
			CreatedAt: time.Now(),
		},
		By: by,
	}
	return mn, nil
}
