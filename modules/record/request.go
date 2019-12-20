package record

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
)

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
func (r *RequestAddRecord) Valid() error {
	if _, err := valid.ValidateStruct(r); err != nil {
		return err
	}
	// 1. 时间检验
	if !r.StartTime.IsZero() && !r.EndTime.IsZero() {
		if r.StartTime.After(r.EndTime) {
			return errors.New("startTime should be before endTime")
		}
	}
	if r.EndTime.After(time.Now()) || r.StartTime.After(time.Now()) {
		return errors.New("cannot add future time")
	}
	// 2. 若公里数不为空时的检验
	if r.StartMileAge != nil && r.EndMileAge != nil && *r.StartMileAge > *r.EndMileAge {
		return errors.New("startMileAge should be less than endMileAge")
	}
	return nil
}
