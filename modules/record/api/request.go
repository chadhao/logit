package api

import (
	"errors"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/record/model"
	userInternals "github.com/chadhao/logit/modules/user/internals"
)

// reqRecords 请求获取记录
type reqRecords struct {
	DriverID            string    `query:"driverID" valid:"required"`
	From                time.Time `query:"from" valid:"required"`
	To                  time.Time `query:"to" valid:"optional"`
	TransportOperatorID string    `query:"transportOperatorID" valid:"-"` // To required
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
	driverID, err := primitive.ObjectIDFromHex(reqR.DriverID)
	if err != nil {
		return nil, err
	}

	// 获取记录
	records, err := model.GetRecords(driverID, reqR.From, reqR.To, false)
	if err != nil {
		return nil, err
	}
	// 获取记录下的笔记以及具体的vehicle信息
	recordIDs := []primitive.ObjectID{}
	vehicleIDs := []primitive.ObjectID{}
	for _, v := range records {
		recordIDs = append(recordIDs, v.ID)
		vehicleIDs = append(vehicleIDs, v.VehicleID)
	}
	notesMap, err := model.GetNotesByRecordIDs(recordIDs)
	if err != nil {
		return nil, err
	}

	vehiclesMap, err := userInternals.GetVehicleMapByIDs(vehicleIDs)
	if err != nil {
		return nil, err
	}

	// 拼装返回
	respRecords := []*respRecord{}
	for _, v := range records {
		respRecords = append(respRecords, &respRecord{
			Record:  v,
			Notes:   notesMap[v.ID],
			Vehicle: vehiclesMap[v.VehicleID],
		})
	}

	return respRecords, nil
}

// reqRecord 请求获取记录
type reqRecord struct {
	ID primitive.ObjectID
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
func (reqRecord *reqRecord) deleteRecord(driverID primitive.ObjectID) error {
	// 获取记录
	r, err := model.GetRecord(reqRecord.ID)
	if err != nil {
		return err
	}
	if r.DriverID != driverID {
		return errors.New("no authorization")
	}

	return r.Delete()
}

// reqAddRecord 添加记录请求结构
type reqAddRecord struct {
	Type          model.Type         `json:"type" valid:"required"`
	Time          time.Time          `json:"time" valid:"required"`
	Duration      string             `json:"duration" valid:"required"`
	StartLocation model.Location     `json:"startLocation" valid:"required"`
	EndLocation   model.Location     `json:"endLocation" valid:"required"`
	VehicleID     primitive.ObjectID `json:"vehicleID" valid:"required"`
	StartMileAge  *float64           `json:"startDistance,omitempty" valid:"-"`
	EndMileAge    *float64           `json:"endDistance,omitempty" valid:"-"`
	ClientTime    *time.Time         `json:"clientTime,omitempty" valid:"-"`
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
	if reqAddR.Time.IsZero() {
		return errors.New("time is required")
	}
	// 如果传入了clientTime,则表示用户当前时间time为用户自己手动输入的时间，则传入不检查(和上条时间对比在下一块儿检查);
	// 如果未传入clientTime,则表示用户当前时间time为用户自己手机的时间，验证是否和标准时间相符，不相符则返回错误。
	if reqAddR.ClientTime == nil {
		if math.Abs(reqAddR.Time.Sub(time.Now()).Seconds()) > 10 {
			return errors.New("time and system time conflict")
		}
	}

	// 2. 若公里数不为空时的检验
	if reqAddR.StartMileAge != nil && reqAddR.EndMileAge != nil && *reqAddR.StartMileAge > *reqAddR.EndMileAge {
		return errors.New("startMileAge should be less than endMileAge")
	}
	return nil
}

// syncValid 上传记录请求结构验证
func (reqAddR *reqAddRecord) syncValid() error {
	if reqAddR.Type != model.WORK && reqAddR.Type != model.REST {
		return errors.New("no match type")
	}
	if reqAddR.ClientTime == nil {
		return errors.New("clientTime is required")
	}
	if _, err := valid.ValidateStruct(reqAddR); err != nil {
		return err
	}
	// 1. 时间检验
	if reqAddR.Time.IsZero() {
		return errors.New("time is required")
	}
	if reqAddR.Time.After(time.Now()) {
		return errors.New("cannot add future time")
	}

	// 2. 若公里数不为空时的检验
	if reqAddR.StartMileAge != nil && reqAddR.EndMileAge != nil && *reqAddR.StartMileAge > *reqAddR.EndMileAge {
		return errors.New("startMileAge should be less than endMileAge")
	}
	return nil
}

// constructToRecord 将reqAddRecord构造为Record
func (reqAddR *reqAddRecord) constructToRecord(driverID primitive.ObjectID) (*model.Record, error) {
	if err := reqAddR.valid(); err != nil {
		return nil, err
	}
	duration, err := time.ParseDuration(reqAddR.Duration)
	if err != nil {
		return nil, err
	}
	r := &model.Record{
		ID:            primitive.NewObjectID(),
		DriverID:      driverID,
		Type:          reqAddR.Type,
		Time:          reqAddR.Time,
		Duration:      duration,
		StartLocation: reqAddR.StartLocation,
		EndLocation:   reqAddR.EndLocation,
		VehicleID:     reqAddR.VehicleID,
		StartMileAge:  reqAddR.StartMileAge,
		EndMileAge:    reqAddR.EndMileAge,
		ClientTime:    reqAddR.ClientTime,
		CreatedAt:     time.Now(),
	}
	return r, nil
}

// constructToSyncRecord 将reqAddRecord构造为上传需要的Record
func (reqAddR *reqAddRecord) constructToSyncRecord(driverID primitive.ObjectID) (*model.Record, error) {
	if err := reqAddR.syncValid(); err != nil {
		return nil, err
	}
	duration, err := time.ParseDuration(reqAddR.Duration)
	if err != nil {
		return nil, err
	}
	r := &model.Record{
		ID:            primitive.NewObjectID(),
		DriverID:      driverID,
		Type:          reqAddR.Type,
		Time:          reqAddR.Time,
		Duration:      duration,
		StartLocation: reqAddR.StartLocation,
		EndLocation:   reqAddR.EndLocation,
		VehicleID:     reqAddR.VehicleID,
		StartMileAge:  reqAddR.StartMileAge,
		EndMileAge:    reqAddR.EndMileAge,
		ClientTime:    reqAddR.ClientTime,
		CreatedAt:     time.Now(),
	}
	return r, nil
}

type reqAddNote struct {
	NoteType            model.NoteType     `json:"noteType" valid:"required"`
	RecordID            primitive.ObjectID `json:"recordID" valid:"required"`
	Comment             string             `json:"comment" valid:"optional"`
	TransportOperatorID primitive.ObjectID `bson:"transportOperatorID" json:"transportOperatorID" valid:"-"`
	StartTime           time.Time          `bson:"startTime" json:"startTime" valid:"-"`
	EndTime             time.Time          `bson:"endTime" json:"endTime" valid:"-"`
	StartLocation       model.Location     `bson:"startLocation" json:"startLocation" valid:"-"`
	EndLocation         model.Location     `bson:"endLocation" json:"endLocation" valid:"-"`
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

// constructToTripNote 将reqAddNote构造为TripNote
func (r *reqAddNote) constructToTripNote() (*model.TripNote, error) {
	if err := r.valid(); err != nil {
		return nil, err
	}
	tn := &model.TripNote{
		Note: model.Note{
			ID:        primitive.NewObjectID(),
			RecordID:  r.RecordID,
			Type:      r.NoteType,
			Comment:   r.Comment,
			CreatedAt: time.Now(),
		},
		TransportOperatorID: r.TransportOperatorID,
		StartTime:           r.StartTime,
		EndTime:             r.EndTime,
		StartLocation:       r.StartLocation,
		EndLocation:         r.EndLocation,
	}
	return tn, nil
}

func (r *reqAddNote) isDriversRecord(driverID primitive.ObjectID) bool {
	rec, err := model.GetRecord(r.RecordID)
	if err != nil {
		return false
	}
	return rec.DriverID == driverID
}
