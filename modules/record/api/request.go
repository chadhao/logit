package api

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/record/model"
)

// RequestRecords 请求获取记录
type RequestRecords struct {
	From time.Time `json:"from" valid:"required"`
	To   time.Time `json:"to" valid:"-"`
}

func (reqR *RequestRecords) valid() error {
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
func (reqR *RequestRecords) getRecords(userID primitive.ObjectID) ([]*ResponseRecord, error) {
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
	respRecords := []*ResponseRecord{}
	notesMap := make(map[primitive.ObjectID][]model.INote)
	for _, v := range notes {
		notesMap[v.GetRecordID()] = append(notesMap[v.GetRecordID()], v)
	}
	for _, v := range records {
		respRecords = append(respRecords, &ResponseRecord{
			Record: v,
			Notes:  notesMap[v.ID],
		})
	}

	return respRecords, nil
}

// RequestRecord 请求获取记录
type RequestRecord struct {
	ID primitive.ObjectID `json:"id" valid:"required"`
}

// getRecord 获取记录
func (reqRecord *RequestRecord) getRecord() (*ResponseRecord, error) {
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
	respRecord := &ResponseRecord{
		Record: *r,
		Notes:  notes,
	}
	// 拼装返回
	return respRecord, nil
}

// RequestAddRecord 添加记录请求结构
type RequestAddRecord struct {
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
func (reqAddR *RequestAddRecord) valid() error {
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

// constructToRecord 将RequestAddRecord构造为Record
func (reqAddR *RequestAddRecord) constructToRecord(userID, vehicleID primitive.ObjectID) (*model.Record, error) {
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

// RequestAddSystemNote 添加系统笔记
type RequestAddSystemNote struct {
	RecordID primitive.ObjectID `json:"recordID" valid:"required"`
	Comment  string             `json:"comment" valid:"required"`
	Type     model.NoteType     `json:"noteType" valid:"required"`
}

// Valid 添加系统笔记验证
func (r *RequestAddSystemNote) valid() error {
	_, err := valid.ValidateStruct(r)
	return err
}

// constructToSystemNote 将RequestAddSystemNote构造为SystemNote
func (r *RequestAddSystemNote) constructToSystemNote() (*model.SystemNote, error) {
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

// RequestAddOtherWorkNote 添加其它笔记
type RequestAddOtherWorkNote struct {
	RecordID primitive.ObjectID `json:"recordID" valid:"required"`
	Comment  string             `json:"comment" valid:"required"`
	Type     model.NoteType     `json:"noteType" valid:"required"`
}

// Valid 添加其它笔记验证
func (r *RequestAddOtherWorkNote) valid() error {
	_, err := valid.ValidateStruct(r)
	return err
}

// constructToOtherWorkNote 将RequestAddOtherWorkNote构造为OtherWorkNote
func (r *RequestAddOtherWorkNote) constructToOtherWorkNote() (*model.OtherWorkNote, error) {
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

// RequestAddModificationNote 添加人为修改笔记
type RequestAddModificationNote struct {
	RecordID primitive.ObjectID `json:"recordID" valid:"required"`
	Comment  string             `json:"comment" valid:"required"`
	Type     model.NoteType     `json:"noteType" valid:"required"`
}

// Valid 添加人为修改笔记
func (r *RequestAddModificationNote) valid() error {
	_, err := valid.ValidateStruct(r)
	return err
}

// constructToModificationNote 将RequestAddModificationNote构造为ModificationNote
func (r *RequestAddModificationNote) constructToModificationNote(by primitive.ObjectID) (*model.ModificationNote, error) {
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
