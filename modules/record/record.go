package record

import (
	"time"

	valid "github.com/asaskevich/govalidator"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// Type 记录类型
	Type string
	// NoteType 笔记类型
	NoteType string
	// HrTime 重要时间点
	HrTime float64
)

const (
	// WORK 工作记录类型
	WORK Type = "work"
	// REST 休息记录类型
	REST Type = "rest"
)
const (
	// SYSTEMNOTE 系统笔记类型
	SYSTEMNOTE NoteType = "system"
	// MODIFICATIONNOTE 人为修改笔记类型
	MODIFICATIONNOTE NoteType = "modification"
	// TRIPNOTE 行程笔记类型
	TRIPNOTE NoteType = "trip"
	// OTHERWORKNOTE 其它笔记类型
	OTHERWORKNOTE NoteType = "others"
)
const (
	// HR0D5 0.5小时
	HR0D5 HrTime = 0.5
	// HR5D5 5.5小时
	HR5D5 HrTime = 5.5
	// HR7 7.5小时
	HR7 HrTime = 7
	// HR10 10小时
	HR10 HrTime = 10
	// HR13 13小时
	HR13 HrTime = 13
	// HR24 24小时
	HR24 HrTime = 24
	// HR70 70小时
	HR70 HrTime = 70
)

func (t HrTime) getHrs() float64 {
	return float64(t)
}

// Location 位置信息
type Location struct {
	Address string     `bson:"address" json:"address"`
	Coors   [2]float64 `bson:"coors,omitempty" json:"coors,omitempty"`
}

// Record 记录
type Record struct {
	ID            primitive.ObjectID `bson:"_id" json:"id" valid:"required"`
	Type          Type               `bson:"type" json:"type" valid:"required"`
	StartTime     time.Time          `bson:"startTime" json:"startTime" valid:"required"`
	EndTime       time.Time          `bson:"endTime,omitempty" json:"endTime,omitempty" valid:"-"`
	StartLocation Location           `bson:"startLocation" json:"startLocation" valid:"required"`
	EndLocation   Location           `bson:"endLocation,omitempty" json:"endLocation,omitempty" valid:"-"`
	VehicleID     primitive.ObjectID `bson:"vehicleID" json:"vehicleID" valid:"required"`
	StartMileAge  *float64           `bson:"startDistance,omitempty" json:"startDistance,omitempty" valid:"-"`
	EndMileAge    *float64           `bson:"endDistance,omitempty" json:"endDistance,omitempty" valid:"-"`
	Notes         []Note             `bson:"notes,omitempty" json:"notes,omitempty" valid:"-"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt" valid:"required"`
	ClientTime    *time.Time         `bson:"clientTime,omitempty" json:"clientTime,omitempty" valid:"-"`
	DeletedAt     *time.Time         `bson:"isDeleted,omitempty" json:"isDeleted,omitempty" valid:"-"`
}

func (r *Record) addNote(note *Note) {
	r.Notes = append(r.Notes, *note)
}

// Add 记录添加
func (r *Record) Add() (err error) {
	if err = r.beforeAdd(); err != nil {
		return
	}
	// 数据库添加记录
	go r.afterAdd()
	return
}
func (r *Record) beforeAdd() error {
	// 1. 验证记录结构是否完整
	return nil
}
func (r *Record) afterAdd() error {
	// 1. 获取上一条记录，将上一条记录信息补充完整; endTime, endLocation不为空时候添加
	// 2. 若此条和上一条的startMileAge都不为空，上一条endMileAge不为空，则为上一条添加endMileAge
	return nil
}
func (r *Record) valid() error {
	if _, err := valid.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

func getLastRecord() (record *Record) {
	return
}
