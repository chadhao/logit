package service

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// UserRegisterInput 用户注册参数
	UserRegisterInput struct {
		Phone    string `json:"phone" valid:"numeric,stringlength(8|11)"`
		Code     string `json:"code" valid:"numeric"`
		Email    string `json:"email" valid:"email"`
		Password string `json:"password" valid:"stringlength(6|32)"`
	}
	// UserRegisterOutput 用户注册返回参数
	UserRegisterOutput struct {
		ID              primitive.ObjectID `json:"id"`
		Phone           string             `json:"phone"`
		Email           string             `json:"email"`
		IsEmailVerified bool               `json:"isEmailVerified"`
		IsDriver        bool               `json:"isDriver"`
		RoleIDs         []int              `json:"roleIDs"`
		CreatedAt       time.Time          `json:"createdAt"`
	}
)

func (u *UserRegisterInput) toUser() *model.User {
	return &model.User{
		ID:        primitive.NewObjectID(),
		Phone:     u.Phone,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: time.Now(),
	}
}

// UserRegister 用户注册
func UserRegister(in *UserRegisterInput) (*UserRegisterOutput, error) {
	// 参数验证
	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}

	// 验证码匹配
	if code, err := model.RedisClient.Get(in.Phone).Result(); err != nil || in.Code != code {
		return nil, errors.New("verification code does not match")
	}

	// 判断用户是否存在
	if model.IsUserExists(model.UserExistsOpt{Phone: in.Phone, Email: in.Email}) {
		return nil, errors.New("user exists")
	}

	// 创建用户
	user := in.toUser()
	if err := user.Create(); err != nil {
		return nil, err
	}

	// 验证码过期处理
	model.RedisClient.ExpireAt(in.Phone, time.Now())

	out := &UserRegisterOutput{
		ID:        user.ID,
		Phone:     user.Phone,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	return out, nil
}

type (
	// UserUpdateInput 用户更新参数
	UserUpdateInput struct {
		UserID   primitive.ObjectID `valid:"required"`
		Password string             `json:"password" valid:"stringlength(6|32),optional"`
		Pin      string             `json:"pin" valid:"stringlength(4|4),optional"`
	}
)

// UserUpdate 用户更新
func UserUpdate(in *UserUpdateInput) error {

	// 查找用户, 验证相关输入, 并且替换相关user信息
	if _, err := valid.ValidateStruct(in); err != nil {
		return err
	}

	user, err := model.FindUser(model.FindUserOpt{UserID: in.UserID})
	if err != nil {
		return err
	}

	if len(in.Password) > 0 {
		user.Password = in.Password
	}
	// 非driver用户不能更新pin
	if len(in.Pin) > 0 {
		if !user.IsDriver {
			return errors.New("only driver need to set pin")
		}
		user.Pin = in.Pin
	}

	// 更新用户
	return user.Update()
}

type (
	// UserExistanceCheckInput 检查用户是否存在参数
	UserExistanceCheckInput struct {
		Phone   string `json:"phone"`
		Email   string `json:"email"`
		License string `json:"license"`
	}
)

// UserExistanceCheck 检查用户是否存在
func UserExistanceCheck(in *UserExistanceCheckInput) map[string]bool {
	result := make(map[string]bool, 0)
	if len(in.Phone) > 0 {
		result["phone"] = model.IsUserExists(model.UserExistsOpt{Phone: in.Phone})
	}
	if len(in.Email) > 0 {
		result["email"] = model.IsUserExists(model.UserExistsOpt{Email: in.Email})
	}
	return result
}

type (
	// UserInfoFindInput 查询用户基础及相关信息参数
	UserInfoFindInput struct {
		UserID primitive.ObjectID
	}
	// UserInfoFindOutput 查询用户基础及相关信息返回参数
	UserInfoFindOutput struct {
		ID              primitive.ObjectID                       `json:"id"`
		Phone           string                                   `json:"phone"`
		Email           string                                   `json:"email"`
		IsEmailVerified bool                                     `json:"isEmailVerified"`
		IsDriver        bool                                     `json:"isDriver"`
		RoleIDs         []int                                    `json:"roleIDs"`
		CreatedAt       time.Time                                `json:"createdAt"`
		Driver          *model.Driver                            `json:"driver,omitempty"`
		Identities      []*model.TransportOperatorIdentityDetail `json:"identities,omitempty"`
	}
)

// Format .
func (r *UserInfoFindOutput) Format(user *model.User, driver *model.Driver, identities []*model.TransportOperatorIdentityDetail) {
	r.ID = user.ID
	r.Phone = user.Phone
	r.Email = user.Email
	r.IsEmailVerified = user.IsEmailVerified
	r.IsDriver = user.IsDriver
	r.RoleIDs = user.RoleIDs
	r.CreatedAt = user.CreatedAt
	if !driver.ID.IsZero() {
		r.Driver = driver
	}
	if len(identities) > 0 {
		r.Identities = identities
	}
}

// UserInfoFind 查询用户基础及相关信息
func UserInfoFind(in *UserInfoFindInput) (*UserInfoFindOutput, error) {
	user, err := model.FindUser(model.FindUserOpt{UserID: in.UserID})
	if err != nil {
		return nil, err
	}

	driver := &model.Driver{}
	if utils.RolesAssert(user.RoleIDs).Is(constant.ROLE_DRIVER) {
		driver, _ = model.FindDriver(model.FindDriverOpt{ID: user.ID})
	}

	tois, _ := model.TransportOperatorIdentityFilter(model.TransportOperatorIdentityFilterOpt{
		UserID: user.ID,
	})

	out := &UserInfoFindOutput{}
	out.Format(user, driver, tois)

	return out, nil
}

type (
	// UserQueryInput 查询用户基础及相关信息参数
	UserQueryInput struct {
		Phone string `json:"phone" valid:"stringlength(4|11)"`
		Email string `json:"email" valid:"stringlength(4|50)"`
	}
	// UserQueryOutput 查询用户基础及相关信息返回参数
	UserQueryOutput struct {
		UserInfo []*UserInfoFindOutput
	}
)

// UserQuery 查询用户群基础及相关信息
func UserQuery(in *UserQueryInput) (*UserQueryOutput, error) {
	// 参数验证
	if _, err := valid.ValidateStruct(in); err != nil {
		return nil, err
	}

	// 查询用户
	users, err := model.FilterUser(model.FilterUserOpt{Phone: in.Phone, Email: in.Email})
	if err != nil {
		return nil, err
	}

	// 查询司机信息
	uids := []primitive.ObjectID{}
	for _, v := range users {
		uids = append(uids, v.ID)
	}
	driversMap := make(map[primitive.ObjectID]*model.Driver)
	drivers, _ := model.FindDrivers(uids)
	for _, driver := range drivers {
		driversMap[driver.ID] = driver
	}

	// 查询角色身份信息
	identitiesMap, _ := model.GetIdentitiesByUserIDs(uids)

	out := &UserQueryOutput{}
	for _, user := range users {
		uinfo := &UserInfoFindOutput{}
		uinfo.Format(user, driversMap[user.ID], identitiesMap[user.ID])
		out.UserInfo = append(out.UserInfo, uinfo)
	}
	return out, nil
}
