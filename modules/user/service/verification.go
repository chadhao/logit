package service

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	msgApi "github.com/chadhao/logit/modules/message/internals"

	"github.com/chadhao/logit/modules/user/constant"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/chadhao/logit/utils"
)

type (
	// CheckVerificationCodeInput 手机验证码验证参数
	CheckVerificationCodeInput struct {
		Phone string `json:"phone" valid:"numeric,stringlength(8|11)"`
		Code  string `json:"code" valid:"numeric"`
	}
)

// CheckVerificationCode 手机验证码验证
func CheckVerificationCode(in *CheckVerificationCodeInput) error {

	if _, err := valid.ValidateStruct(in); err != nil {
		return err
	}

	code, err := model.RedisClient.Get(in.Phone).Result()
	if err != nil {
		return err
	}
	if in.Code != code {
		return errors.New("verification code does not match")
	}

	return nil
}

type (
	// EmailVerifyInput 邮箱验证参数
	EmailVerifyInput struct {
		Email string `query:"email" valid:"email"`
		Token string `query:"token" valid:"required"`
	}
	// EmailVerifyOutput 邮箱验证返回参数
	EmailVerifyOutput struct {
		HTML string
	}
)

// EmailVerify 邮箱验证
func EmailVerify(in *EmailVerifyInput) (out *EmailVerifyOutput, err error) {
	out = &EmailVerifyOutput{
		HTML: "<h1>Hi there,</h1><p>Your email has been verified!</p>",
	}

	if _, err = valid.ValidateStruct(in); err != nil {
		out.HTML = "<h1>Bad request</h1><p>" + err.Error() + "</p>"
		return
	}

	// 检查token是否准确
	var token string
	token, err = model.RedisClient.Get(in.Email).Result()
	if err != nil {
		out.HTML = "<h1>Bad request</h1><p>" + err.Error() + "</p>"
		return
	}
	if in.Token != token {
		err = errors.New("token does not match")
		out.HTML = "<h1>Bad request</h1><p>" + err.Error() + "</p>"
		return
	}

	// 更新User相关信息
	user, err := model.FindUser(model.FindUserOpt{Email: in.Email})
	if err != nil {
		out.HTML = "<h1>Bad request</h1><p>" + err.Error() + "</p>"
		return
	}
	user.IsEmailVerified = true
	if err = user.Update(); err != nil {
		out.HTML = "<h1>Bad request</h1><p>" + err.Error() + "</p>"
		return
	}

	model.RedisClient.ExpireAt(in.Email, time.Now())

	// Token过期处理
	return
}

type (
	// SendVerificationInput 获取验证码参数
	SendVerificationInput struct {
		Phone string `json:"phone" valid:"numeric,stringlength(8|11),optional"`
		Email string `json:"email" valid:"email,optional"`
	}
)

func (v *SendVerificationInput) txtSent() (string, error) {
	code := utils.GetRandomCode(6)
	msg := "[LOGIT] Your verification code is: " + code
	return code, msgApi.SendTxt(msgApi.TxtRequest{Number: v.Phone, Message: msg})
}

func (v *SendVerificationInput) emailSent() (string, error) {
	code := utils.GetMD5Hash(v.Email)
	email := msgApi.EmailRequest{
		Sender:     constant.EMAIL_SENDER,
		Recipients: []string{v.Email},
		Subject:    "Logit Verification Email",
		HTMLBody: "<h1>Logit Verification Email</h1><p>Please click " +
			"<a href='https://dev.ssh.logit.co.nz/email/verification?email=" + v.Email + "&token=" + code + "'>here</a>" +
			" to active email.</p>",
		CharSet: "UTF-8",
	}
	return code, msgApi.SendEmail(email)
}

// SendVerification 获取验证码
func SendVerification(in *SendVerificationInput) (err error) {

	if _, err := valid.ValidateStruct(in); err != nil {
		return err
	}

	// 生成code,并保存至redis
	var code, redisKey, durationStr string

	// 发送至电话或者邮箱
	switch {
	case len(in.Phone) > 0:
		redisKey = in.Phone
		durationStr = "10m"
		if code, err = in.txtSent(); err != nil {
			return err
		}
	case valid.IsEmail(in.Email):
		redisKey = in.Email
		durationStr = "12h"
		if code, err = in.emailSent(); err != nil {
			return err
		}
	default:
		return errors.New("phone number or email is requried")
	}

	duration, _ := time.ParseDuration(durationStr)
	model.RedisClient.Set(redisKey, code, duration)

	return nil
}
