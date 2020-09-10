package service

import (
	"errors"
	"math"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/record/model"
	userInternals "github.com/chadhao/logit/modules/user/internals"
	userModel "github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// CreateRecordInput 新建记录参数
	CreateRecordInput struct {
		DriverID      primitive.ObjectID `json:"driverID" valid:"required"`
		Type          model.Type         `json:"type" valid:"required"`
		Time          time.Time          `json:"time" valid:"required"`
		Duration      time.Duration      `json:"duration" valid:"required"`
		StartLocation model.Location     `json:"startLocation" valid:"required"`
		EndLocation   model.Location     `json:"endLocation" valid:"required"`
		VehicleID     primitive.ObjectID `json:"vehicleID" valid:"required"`
		StartMileAge  *float64           `json:"startDistance,omitempty" valid:"-"`
		EndMileAge    *float64           `json:"endDistance,omitempty" valid:"-"`
		ClientTime    *time.Time         `json:"clientTime,omitempty" valid:"-"`
	}
	// CreateRecordOutput 新建记录返回参数
	CreateRecordOutput struct {
		*model.Record `json:",inline"`
	}
)

// valid 添加记录请求结构验证
func (r *CreateRecordInput) valid() error {
	if r.Type != model.WORK && r.Type != model.REST {
		return errors.New("no match type")
	}
	if r.VehicleID.IsZero() {
		return errors.New("vehicleID is required")
	}
	if _, err := valid.ValidateStruct(r); err != nil {
		return err
	}

	// 1. 时间检验
	if r.Time.IsZero() {
		return errors.New("time is required")
	}
	// 如果传入了clientTime,则表示用户当前时间time为用户自己手动输入的时间，则传入不检查(和上条时间对比在下一块儿检查);
	// 如果未传入clientTime,则表示用户当前时间time为用户自己手机的时间，验证是否和标准时间相符，不相符则返回错误。
	if r.ClientTime == nil {
		if math.Abs(r.Time.Sub(time.Now()).Seconds()) > 10 {
			return errors.New("time and system time conflict")
		}
	}

	// 2. 若公里数不为空时的检验
	if r.StartMileAge != nil && r.EndMileAge != nil && *r.StartMileAge > *r.EndMileAge {
		return errors.New("startMileAge should be less than endMileAge")
	}

	// 3. 检查endLocation
	if err := r.EndLocation.FillFull(); err != nil {
		return err
	}
	return nil
}

// syncValid 上传记录请求结构验证
func (r *CreateRecordInput) syncValid() error {
	if r.Type != model.WORK && r.Type != model.REST {
		return errors.New("no match type")
	}
	if r.VehicleID.IsZero() {
		return errors.New("vehicleID is required")
	}
	if r.ClientTime == nil {
		return errors.New("clientTime is required")
	}
	if _, err := valid.ValidateStruct(r); err != nil {
		return err
	}
	// 1. 时间检验
	if r.Time.IsZero() {
		return errors.New("time is required")
	}
	if r.Time.After(time.Now()) {
		return errors.New("cannot add future time")
	}

	// 2. 若公里数不为空时的检验
	if r.StartMileAge != nil && r.EndMileAge != nil && *r.StartMileAge > *r.EndMileAge {
		return errors.New("startMileAge should be less than endMileAge")
	}
	// 3. 检查endLocation
	if err := r.EndLocation.FillFull(); err != nil {
		return err
	}
	return nil
}

// validWithLastRec 添加记录与上一条记录冲突验证
func (r *CreateRecordInput) validWithLastRec(lastRec *model.Record) error {
	if lastRec.Type == r.Type {
		return errors.New("work type conflict with last record")
	}
	if lastRec.Time.After(r.Time) {
		return errors.New("time conflict with last record")
	}
	if !r.StartLocation.Equal(&lastRec.EndLocation) {
		return errors.New("location not match")
	}
	if r.StartMileAge != nil && !utils.AlmostEqual(*r.StartMileAge, *lastRec.EndMileAge) {
		return errors.New("mileage not match")
	}
	if math.Abs(lastRec.Time.Add(r.Duration).Sub(r.Time).Seconds()) > 10 {
		return errors.New("time and duration not match")
	}
	return nil
}

// CreateRecord 新建记录
func CreateRecord(in *CreateRecordInput) (*CreateRecordOutput, error) {
	// 获取最近的一条记录，如果错误原因并非无上一条记录不存在，则返回错误
	lastRec, err := model.GetLastestRecord(in.DriverID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	// 新纪录验证
	if err := in.valid(); err != nil {
		return nil, err
	}

	// 若上一条记录存在，则新纪录与上一条记录比对验证
	if (model.Record{}) != *lastRec {
		if err := in.validWithLastRec(lastRec); err != nil {
			return nil, err
		}
	}

	// 添加新记录
	active := true
	newRecord := &model.Record{
		ID:            primitive.NewObjectID(),
		DriverID:      in.DriverID,
		Type:          in.Type,
		Time:          in.Time,
		Duration:      in.Duration,
		StartLocation: in.StartLocation,
		EndLocation:   in.EndLocation,
		VehicleID:     in.VehicleID,
		StartMileAge:  in.StartMileAge,
		EndMileAge:    in.EndMileAge,
		ClientTime:    in.ClientTime,
		CreatedAt:     time.Now(),
		Active:        &active,
	}

	if err := newRecord.Add(lastRec); err != nil {
		return nil, err
	}

	return &CreateRecordOutput{newRecord}, nil
}

// SyncRecords 同步多条记录
func SyncRecords(in []*CreateRecordInput) ([]*CreateRecordOutput, error) {
	// 获取最近的一条记录，如果错误原因并非无上一条记录不存在，则返回错误
	lastRec, err := model.GetLastestRecord(in[0].DriverID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	lastRecCopy := *lastRec
	newRecords := model.Records{}
	l := len(in)
	if l < 1 {
		return nil, errors.New("at least one record id required")
	}
	for i := 0; i < l; i++ {
		// 新纪录验证
		if err = in[i].valid(); err != nil {
			return nil, err
		}
		// 若上一条记录存在，则新纪录与上一条记录比对验证
		if (model.Record{}) != *lastRec {
			if err = in[i].validWithLastRec(lastRec); err != nil {
				return nil, err
			}
			newRecord := &model.Record{
				ID:            primitive.NewObjectID(),
				DriverID:      in[i].DriverID,
				Type:          in[i].Type,
				Time:          in[i].Time,
				Duration:      in[i].Duration,
				StartLocation: in[i].StartLocation,
				EndLocation:   in[i].EndLocation,
				VehicleID:     in[i].VehicleID,
				StartMileAge:  in[i].StartMileAge,
				EndMileAge:    in[i].EndMileAge,
				ClientTime:    in[i].ClientTime,
				CreatedAt:     time.Now(),
			}
			newRecords = append(newRecords, newRecord)
			lastRec = newRecord
		}
	}
	// 为最后一条record添加active
	active := true
	newRecords[l-1].Active = &active

	if err := newRecords.SyncAdd(&lastRecCopy); err != nil {
		return nil, err
	}

	out := []*CreateRecordOutput{}
	for i := 0; i < l; i++ {
		out = append(out, &CreateRecordOutput{newRecords[i]})
	}
	return out, nil
}

type (
	// GetRecordsInput 请求获取记录参数
	GetRecordsInput struct {
		DriverID            string    `query:"driverID" valid:"required"`
		From                time.Time `query:"from" valid:"required"`
		To                  time.Time `query:"to" valid:"optional"`
		TransportOperatorID string    `query:"transportOperatorID" valid:"-"` // To required
	}
	// GetRecordOut 获取记录及其notes返回
	GetRecordOut struct {
		model.Record `json:",inline"`
		Notes        model.DifNotes     `json:"notes,omitempty"`
		Vehicle      *userModel.Vehicle `json:"vehicle"`
	}
)

// GetRecords 请求获取记录
func GetRecords(in *GetRecordsInput) ([]*GetRecordOut, error) {
	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}
	driverID, err := primitive.ObjectIDFromHex(in.DriverID)
	if err != nil {
		return nil, err
	}

	// 获取记录
	records, err := model.GetRecords(driverID, model.GetRecordsOpt{From: in.From, To: in.To, GetDeleted: false})

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
	out := []*GetRecordOut{}
	for _, v := range records {
		out = append(out, &GetRecordOut{
			Record:  *v,
			Notes:   notesMap[v.ID],
			Vehicle: vehiclesMap[v.VehicleID],
		})
	}

	return out, nil
}

// GetLastestRecord 获取最近的一条记录及其notes
func GetLastestRecord(driverID primitive.ObjectID) (*GetRecordOut, error) {
	lastRec, err := model.GetLastestRecord(driverID)
	if err != nil {
		return nil, err
	}

	notesMap, err := model.GetNotesByRecordIDs([]primitive.ObjectID{lastRec.ID})
	if err != nil {
		return nil, err
	}

	out := &GetRecordOut{
		Record: *lastRec,
		Notes:  notesMap[lastRec.ID],
	}
	return out, nil
}

// DeleteLatestRecordInput 删除最近的一条记录参数
type DeleteLatestRecordInput struct {
	RecordID primitive.ObjectID
	DriverID primitive.ObjectID
}

// DeleteLatestRecord 删除最近的一条记录
func DeleteLatestRecord(in *DeleteLatestRecordInput) error {
	// 获取记录
	r, err := model.GetRecord(in.RecordID)
	if err != nil {
		return err
	}
	// 验证记录
	if r.DriverID != in.DriverID {
		return errors.New("driver has no authorization")
	}
	if r.DeletedAt != nil {
		return errors.New("record has already been deleted")
	}

	lastRec, err := model.GetLastestRecord(r.DriverID)
	if err != mongo.ErrNoDocuments {
		if err != nil {
			return err
		}
		if lastRec.ID != r.ID {
			return errors.New("record is not the lastest one")
		}
	}
	// 删除记录
	return r.Delete(lastRec)
}
