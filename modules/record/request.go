package record

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	valid "github.com/asaskevich/govalidator"
)

// GetLatest7DaysRecords 获取前7天时间范围内的记录
func GetLatest7DaysRecords(userID primitive.ObjectID) ([]Record, error) {
	to := time.Now()
	from := to.AddDate(0, 0, 7)
	return getRecords(userID, from, to, false)
}

// GetRecord 获取记录
func GetRecord(id primitive.ObjectID) (*Record, error) {
	return getRecord(id)
}

// RequestAddRecord 添加记录请求结构
type RequestAddRecord struct {
	Type          Type       `json:"type" valid:"required"`
	StartTime     time.Time  `json:"startTime,omitempty" valid:"-"`
	EndTime       time.Time  `json:"endTime,omitempty" valid:"-"`
	StartLocation Location   `json:"startLocation" valid:"required"`
	EndLocation   Location   `json:"endLocation,omitempty" valid:"-"`
	StartMileAge  *float64   `json:"startDistance,omitempty" valid:"-"`
	EndMileAge    *float64   `json:"endDistance,omitempty" valid:"-"`
	ClientTime    *time.Time `bson:"clientTime,omitempty" json:"clientTime,omitempty" valid:"-"`
}

// Valid 添加记录请求结构验证
func (rar *RequestAddRecord) Valid() error {
	if _, err := valid.ValidateStruct(rar); err != nil {
		return err
	}
	// 1. 时间检验
	if !rar.StartTime.IsZero() && !rar.EndTime.IsZero() {
		if rar.StartTime.After(rar.EndTime) {
			return errors.New("startTime should be before endTime")
		}
	}
	if rar.EndTime.After(time.Now()) || rar.StartTime.After(time.Now()) {
		return errors.New("cannot add future time")
	}
	// 2. 若公里数不为空时的检验
	if rar.StartMileAge != nil && rar.EndMileAge != nil && *rar.StartMileAge > *rar.EndMileAge {
		return errors.New("startMileAge should be less than endMileAge")
	}
	return nil
}

// ConstructToRecord 将RequestAddRecord构造为Record
func (rar *RequestAddRecord) ConstructToRecord(userID primitive.ObjectID) *Record {
	now := time.Now()
	r := &Record{
		ID:            primitive.NewObjectID(),
		UserID:        userID,
		Type:          rar.Type,
		StartTime:     rar.StartTime,
		EndTime:       rar.EndTime,
		StartLocation: rar.StartLocation,
		EndLocation:   rar.EndLocation,
		// 获取vehivleID
		// VehicleID: user.GetVehivleID(),
		StartMileAge: rar.StartMileAge,
		EndMileAge:   rar.EndMileAge,
		ClientTime:   rar.ClientTime,
		CreatedAt:    now,
	}
	if r.StartTime.IsZero() {
		r.StartTime = now
	}
	return r
}
