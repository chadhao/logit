package service

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// DriverRegisterInput 司机注册参数
	DriverRegisterInput struct {
		Conf          config.Config
		UserID        primitive.ObjectID
		LicenseNumber string    `json:"licenseNumber" valid:"stringlength(5|8)"`
		DateOfBirth   time.Time `json:"dateOfBirth" valid:"required"`
		Firstnames    string    `json:"firstnames" valid:"required"`
		Surname       string    `json:"surname" valid:"required"`
		Pin           string    `json:"pin" valid:"stringlength(4|4)"`
	}
)

func (d *DriverRegisterInput) toDriver() *model.Driver {
	return &model.Driver{
		ID:            d.UserID,
		LicenseNumber: d.LicenseNumber,
		DateOfBirth:   d.DateOfBirth,
		Firstnames:    d.Firstnames,
		Surname:       d.Surname,
		CreatedAt:     time.Now(),
	}
}

// DriverRegister 司机注册
func DriverRegister(in *DriverRegisterInput) (*IssueTokenOutput, error) {
	// 参数验证
	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}

	// 查询用户
	user, err := model.FindUser(model.FindUserOpt{UserID: in.UserID})
	if err != nil {
		return nil, err
	}
	if user.IsDriver {
		return nil, errors.New("is driver already")
	}

	// 司机身份注册
	driver := in.toDriver()
	if model.IsDriverExists(model.DriverExistsOpt{LicenseNumber: driver.LicenseNumber}) {
		return nil, errors.New("licenseNumber has been used")
	}

	// 创建司机身份，并且更新用户信息
	user.IsDriver = true
	user.RoleIDs = append(user.RoleIDs, constant.ROLE_DRIVER)
	user.Pin = in.Pin // driver用户需要设置pin
	if err := driver.Create(user); err != nil {
		return nil, err
	}

	// 发放新的token
	issueTokenOutput, err := IssueToken(&IssueTokenInput{UserID: user.ID, RoleIDs: user.RoleIDs, Conf: in.Conf})
	if err != nil {
		return nil, err
	}
	return issueTokenOutput, nil
}

// DriverPinCheckInput 司机Pin验证参数
type DriverPinCheckInput struct {
	UserID primitive.ObjectID
	Pin    string `json:"pin" valid:"stringlength(4|4)"`
}

// DriverPinCheck 司机Pin验证
func DriverPinCheck(in *DriverPinCheckInput) error {
	// 参数验证
	if _, err := valid.ValidateStruct(in); err != nil {
		return err
	}

	// 查询用户
	user, err := model.FindUser(model.FindUserOpt{UserID: in.UserID})
	if err != nil {
		return err
	}

	if user.Pin != in.Pin {
		return errors.New("not match")
	}
	return nil
}

type (
	// DriversFindByTOInput 查询TO下的司机信息参数
	DriversFindByTOInput struct {
		OperatorID          primitive.ObjectID
		TransportOperatorID primitive.ObjectID
	}
	// DriversFindByTOOut 查询TO下的司机信息返回参数
	DriversFindByTOOut struct {
		Drivers []*model.Driver
	}
)

// DriversFindByTO 查询TO下的司机信息
func DriversFindByTO(in *DriversFindByTOInput) (*DriversFindByTOOut, error) {
	// 检查操作者是否有TO权限
	operator, err := model.FindUser(model.FindUserOpt{UserID: in.OperatorID})
	if err != nil {
		return nil, err
	}
	// 如果operator不是admin以上，则如果不具有该TO的admin以上权限，返回错误
	if !utils.RolesAssert(operator.RoleIDs).Are([]int{constant.ROLE_ADMIN, constant.ROLE_SUPER}) {
		if !model.IsTransportOperatorIdentityExists(model.TransportOperatorIdentityExists{
			UserID:              in.OperatorID,
			TransportOperatorID: in.TransportOperatorID,
			Identity:            []model.TOIdentity{model.TO_ADMIN, model.TO_SUPER},
		}) {
			return nil, errors.New("operator has no authorization")
		}
	}

	// 查询用户
	tois, err := model.TransportOperatorIdentityFilter(model.TransportOperatorIdentityFilterOpt{
		TransportOperatorID: in.TransportOperatorID,
		Identity:            model.TO_DRIVER,
	})
	if err != nil {
		return nil, err
	}

	driverIDs := []primitive.ObjectID{}
	for _, v := range tois {
		driverIDs = append(driverIDs, v.UserID)
	}
	drivers, err := model.FindDrivers(driverIDs)

	return &DriversFindByTOOut{drivers}, nil
}
