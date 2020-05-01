package request

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/config"
	msgApi "github.com/chadhao/logit/modules/message/internals"

	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/utils"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	RefreshTokenRequest struct {
		Token string `json:"token"`
	}
	LoginRequest struct {
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		License  string `json:"license"`
		Password string `json:"password"`
	}
	ExistanceRequest struct {
		Phone   string `json:"phone"`
		Email   string `json:"email"`
		License string `json:"license"`
	}
	VerificationRequest struct {
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	EmailVerifyRequest struct {
		Email string `query:"email" valid:"email"`
		Token string `query:"token" valid:"required"`
	}
	ForgetPasswordRequest struct {
		Phone    string `json:"phone" valid:"numeric,stringlength(8|11),optional"`
		Email    string `json:"email" valid:"email,optional"`
		Token    string `json:"token" valid:"required"`
		Password string `json:"password" valid:"stringlength(6|32)"`
	}
)

func (r *RefreshTokenRequest) Validate(c config.Config) (*model.User, error) {
	u := model.User{}

	key, _ := c.Get("system.jwt.refresh.key")
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}
	token, err := jwt.Parse(r.Token, keyFunc)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	userID, err := primitive.ObjectIDFromHex(claims["sub"].(string))
	if err != nil {
		return nil, err
	}
	u.ID = userID

	return &u, nil
}

func (r *LoginRequest) PasswordLogin() (*model.User, error) {
	u := model.User{}

	if len(r.Phone) > 0 || len(r.Email) > 0 {
		u.Phone = r.Phone
		u.Email = r.Email
		u.Password = r.Password
		if err := u.PasswordLogin(); err != nil {
			return nil, err
		}
	} else {
		d := model.Driver{
			LicenseNumber: r.License,
		}
		if err := d.Find(); err != nil {
			return nil, errors.New("user not found")
		}
		u.ID = d.ID
		u.Password = r.Password
		if err := u.PasswordLogin(); err != nil {
			return nil, err
		}
	}

	return &u, nil
}

func (e *EmailVerifyRequest) Verify() (*model.User, error) {
	if _, err := valid.ValidateStruct(e); err != nil {
		return nil, err
	}

	red := model.Redis{Key: e.Email}
	if token, err := red.Get(); err != nil || e.Token != token {
		return nil, errors.New("token does not match")
	}

	u := &model.User{
		Email: e.Email,
	}

	if err := u.Find(); err != nil {
		return nil, err
	}

	u.IsEmailVerified = true
	if err := u.Update(); err != nil {
		return nil, err
	}

	red.Expire()

	return u, nil
}

func (r *ExistanceRequest) Check() map[string]bool {
	result := make(map[string]bool, 0)
	if len(r.Phone) > 0 {
		u := model.User{
			Phone: r.Phone,
		}
		result["phone"] = u.Exists()
	}
	if len(r.Email) > 0 {
		u := model.User{
			Email: r.Email,
		}
		result["email"] = u.Exists()
	}
	if len(r.License) > 0 {
		d := model.Driver{
			LicenseNumber: r.License,
		}
		result["licemnce"] = d.Exists()
	}
	return result
}

func (r *VerificationRequest) Send() (err error) {
	// 生成code,并保存至redis
	var code, redisKey, durationStr string
	// 发送至电话或者邮箱
	switch {
	case len(r.Phone) > 0 && valid.IsNumeric(r.Phone):
		redisKey = r.Phone
		durationStr = "10m"
		if code, err = r.txtSent(); err != nil {
			return err
		}
	case valid.IsEmail(r.Email):
		redisKey = r.Email
		durationStr = "12h"
		if code, err = r.emailSent(); err != nil {
			return err
		}
	default:
		return errors.New("phone number or email is requried")
	}

	duration, _ := time.ParseDuration(durationStr)
	red := model.Redis{
		Key:            redisKey,
		ExpireDuration: duration,
	}
	red.Set(code)

	return nil
}

func (r *VerificationRequest) txtSent() (string, error) {
	code := utils.GetRandomCode(6)
	msg := "[LOGIT] Your verification code is: " + code
	return code, msgApi.SendTxt(msgApi.TxtRequest{Number: r.Phone, Message: msg})
}

func (r *VerificationRequest) emailSent() (string, error) {
	code := utils.GetMD5Hash(r.Email)
	email := msgApi.EmailRequest{
		Sender:     constant.EMAIL_SENDER,
		Recipients: []string{r.Email},
		Subject:    "Logit Verification Email",
		HTMLBody: "<h1>Logit Verification Email</h1><p>Please click " +
			"<a href='http://dev.logit.co.nz/email/verification?email=" + r.Email + "&token=" + code + "'>here</a>" +
			" to active email.</p>",
		CharSet: "UTF-8",
	}
	return code, msgApi.SendEmail(email)
}

func (r *ForgetPasswordRequest) Verify() (err error) {
	if _, err := valid.ValidateStruct(r); err != nil {
		return err
	}
	var redisKey string
	if len(r.Phone) > 0 {
		redisKey = r.Phone
	} else if len(r.Email) > 0 {
		redisKey = r.Email
	} else {
		return errors.New("phone number or email is required")
	}
	red := model.Redis{Key: redisKey}
	if token, err := red.Get(); err != nil || r.Token != token {
		return errors.New("token does not match")
	}

	red.Expire()
	return nil
}
