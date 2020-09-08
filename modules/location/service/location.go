package service

import (
	"errors"
	"fmt"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/location/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// CreateDrivingLocInput 创建司机在行驶过程中的位置信息
	CreateDrivingLocInput struct {
		DriverID  primitive.ObjectID `json:"driverID" valid:"required"`
		Lat       float64            `json:"lat" valid:"required"`
		Lng       float64            `json:"lng" valid:"required"`
		CreatedAt time.Time          `json:"createdAt" valid:"-"`
	}
	// CreateDrivingLocOutput 创建司机在行驶过程中的位置信息返回参数
	CreateDrivingLocOutput struct {
		*model.DrivingLoc `json:",inline"`
	}
)

// CreateDrivingLoc 创建司机在行驶过程中的位置信息
func CreateDrivingLoc(in *CreateDrivingLocInput) (*CreateDrivingLocOutput, error) {

	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}
	if in.DriverID.IsZero() {
		return nil, errors.New("driverID is required")
	}
	if !valid.IsLongitude(fmt.Sprintf("%f", in.Lng)) || !valid.IsLatitude(fmt.Sprintf("%f", in.Lat)) {
		return nil, errors.New("valid coors is required")
	}
	if in.CreatedAt.IsZero() {
		in.CreatedAt = time.Now()
	}

	drivinLoc := &model.DrivingLoc{
		ID:        primitive.NewObjectID(),
		DriverID:  in.DriverID,
		CreatedAt: in.CreatedAt,
		Coors:     model.Coors{Lat: in.Lat, Lng: in.Lng},
	}
	if err := drivinLoc.Create(); err != nil {
		return nil, err
	}

	return &CreateDrivingLocOutput{drivinLoc}, nil
}

type (
	// FindDrivingLocsInput 获取司机行驶的位置信息参数
	FindDrivingLocsInput struct {
		DriverID primitive.ObjectID `json:"driverID" query:"driverID" valid:"required"`
		From     time.Time          `json:"coors" query:"from" valid:"required"`
		To       time.Time          `json:"createdAt" query:"to" valid:"optional"`
	}
	// FindDrivingLocsOutput 获取司机行驶的位置信息返回参数
	FindDrivingLocsOutput struct {
		DrivingLocs []*model.DrivingLoc `json:"drivingLocs"`
	}
)

// FindDrivingLocs 获取司机行驶的位置信息
func FindDrivingLocs(in *FindDrivingLocsInput) (*FindDrivingLocsOutput, error) {
	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}
	if in.DriverID.IsZero() {
		return nil, errors.New("driverID is required")
	}
	if in.To.IsZero() {
		in.To = time.Now()
	}

	out, err := model.GetDrivingLocs(in.DriverID, in.From, in.To)
	if err != nil {
		return nil, err
	}
	return &FindDrivingLocsOutput{out}, nil
}
