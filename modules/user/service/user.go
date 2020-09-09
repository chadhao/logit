package service

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/user/model"
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
