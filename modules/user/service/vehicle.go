package service

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// VehicleCreateInput 创建车辆信息参数
	VehicleCreateInput struct {
		DriverID     primitive.ObjectID `json:"driverID" valid:"required"`
		Registration string             `json:"registration" valid:"stringlength(5|8)"`
		IsDiesel     bool               `json:"isDiesel" valid:"-"`
	}
	// VehicleCreateOutput 创建车辆信息返回参数
	VehicleCreateOutput struct {
		*model.Vehicle
	}
)

// VehicleCreate 创建车辆信息
func VehicleCreate(in *VehicleCreateInput) (*VehicleCreateOutput, error) {

	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}

	vehicle := &model.Vehicle{
		ID:           primitive.NewObjectID(),
		DriverID:     in.DriverID,
		Registration: in.Registration,
		IsDiesel:     in.IsDiesel,
		CreatedAt:    time.Now(),
	}

	if model.IsVehicleExists(model.VehicleExistsOpt{DriverID: in.DriverID, Registration: in.Registration}) {
		return nil, errors.New("vehivle exists")
	}

	if err := vehicle.Create(); err != nil {
		return nil, err
	}

	return &VehicleCreateOutput{vehicle}, nil
}

type (
	// VehicleGetInput 获取车辆信息参数
	VehicleGetInput struct {
		VehicleID primitive.ObjectID `valid:"required"`
		UserID    primitive.ObjectID `valid:"required"`
	}
	// VehicleGetOutput 获取车辆信息返回参数
	VehicleGetOutput struct {
		*model.Vehicle
	}
)

// VehicleGet 获取车辆信息
func VehicleGet(in *VehicleGetInput) (*VehicleGetOutput, error) {

	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}

	vehicle, err := model.FindVehicle(in.VehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle.DriverID != in.UserID {
		return nil, errors.New("no authorization")
	}

	return &VehicleGetOutput{vehicle}, nil
}

type (
	// VehiclesGetInput 获取车辆信息参数
	VehiclesGetInput struct {
		DriverID primitive.ObjectID `valid:"required"`
	}
	// VehiclesGetOutput 获取车辆信息返回参数
	VehiclesGetOutput struct {
		Vehicles []*model.Vehicle `json:"vehicles"`
	}
)

// VehiclesGet 获取车辆信息
func VehiclesGet(in *VehiclesGetInput) (*VehiclesGetOutput, error) {

	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}

	vehicles, err := model.FindVehicles(model.FindVehiclesOpt{DriverID: in.DriverID})
	if err != nil {
		return nil, err
	}

	return &VehiclesGetOutput{vehicles}, nil
}

type (
	// VehicleDeleteInput 删除车辆信息参数
	VehicleDeleteInput struct {
		VehicleID primitive.ObjectID `json:"id" valid:"required"`
		UserID    primitive.ObjectID `valid:"required"`
	}
)

// VehicleDelete 删除车辆信息
func VehicleDelete(in *VehicleDeleteInput) error {

	if _, err := valid.ValidateStruct(in); err != nil {
		return err
	}

	vehicle, err := model.FindVehicle(in.VehicleID)
	if err != nil {
		return err
	}
	if vehicle.DriverID != in.UserID {
		return errors.New("no authorization")
	}

	return vehicle.Delete()
}
