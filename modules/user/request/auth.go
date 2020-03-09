package request

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/config"
	msgInternal "github.com/chadhao/logit/modules/message/internal"
	msgModel "github.com/chadhao/logit/modules/message/model"

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
		Licence  string `json:"licence"`
		Password string `json:"password"`
	}
	UserRegRequest struct {
		Phone    string `json:"phone"`
		Code     string `json:"code"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	DriverRegRequest struct {
		Id            primitive.ObjectID `json:"id"`
		LicenceNumber string             `json:"licenceNumber"`
		DateOfBirth   time.Time          `json:"dateOfBirth"`
		Firstnames    string             `json:"firstnames"`
		Surname       string             `json:"surname"`
	}
	TransportOperatorRegRequest struct {
		Id            primitive.ObjectID `json:"id"`
		LicenceNumber string             `json:"licenceNumber"`
		Name          string             `json:"name"`
	}
	ExistanceRequest struct {
		Phone   string `json:"phone"`
		Email   string `json:"email"`
		Licence string `json:"licence"`
	}
	VerificationRequest struct {
		Phone string `json:"phone"`
		Email string `json:"email"`
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
	userId, err := primitive.ObjectIDFromHex(claims["sub"].(string))
	if err != nil {
		return nil, err
	}
	u.Id = userId

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
			LicenceNumber: r.Licence,
		}
		if err := d.Find(); err != nil {
			return nil, err
		}
		u.Id = d.Id
		u.Password = r.Password
		if err := u.PasswordLogin(); err != nil {
			return nil, err
		}
	}

	return &u, nil
}

func (r *UserRegRequest) Reg() (*model.User, error) {
	// Should add Request content validation here

	red := model.Redis{Key: r.Phone}
	if code, err := red.Get(); err != nil || r.Code != code {
		return nil, errors.New("verification code does not match")
	}

	u := model.User{
		Phone:    r.Phone,
		Email:    r.Email,
		Password: r.Password,
	}

	if err := u.Create(); err != nil {
		return nil, err
	}

	red.Expire()

	return &u, nil
}

func (r *DriverRegRequest) Reg() (*model.Driver, error) {
	// Should add Request content validation here
	d := model.Driver{
		Id:            r.Id,
		LicenceNumber: r.LicenceNumber,
		DateOfBirth:   r.DateOfBirth,
		Firstnames:    r.Firstnames,
		Surname:       r.Surname,
	}

	if err := d.Create(); err != nil {
		return nil, err
	}

	return &d, nil
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
	if len(r.Licence) > 0 {
		d := model.Driver{
			LicenceNumber: r.Licence,
		}
		result["licemnce"] = d.Exists()
	}
	return result
}

func (r *VerificationRequest) Send() (err error) {
	// 生成code,并保存至redis
	var code, redisKey string

	// 发送至电话或者邮箱
	switch {
	case valid.IsNumeric(r.Phone):
		redisKey = r.Phone
		if code, err = r.txtSent(); err != nil {
			return err
		}
	case valid.IsEmail(r.Email):
		redisKey = r.Email
		if code, err = r.emailSent(); err != nil {
			return err
		}
	default:
		return errors.New("phone number or email is requried")
	}

	duration, _ := time.ParseDuration("5m")
	red := model.Redis{
		Key:            redisKey,
		ExpireDuration: duration,
	}
	red.Set(code)

	return nil
}

func (r *VerificationRequest) txtSent() (string, error) {
	code := utils.GetRandomCode(6)
	msg := "[Logit]Your verification code is: " + code
	return code, msgInternal.SendTxt(msgModel.Txt{Number: r.Phone, Message: msg})
}

func (r *VerificationRequest) emailSent() (string, error) {
	code := utils.GetMD5Hash(r.Email)
	email := msgModel.Email{
		Sender:     "sender@logit.co.nz",
		Recipients: []string{r.Email},
		Subject:    "Logit Verification Email",
		HTMLBody: "<h1>Logit Verification Email</h1><p>Please click " +
			"<a href='https://logit.co.nz/email/verification" + code + "'>here</a>" +
			"to active email.</p>",
		CharSet: "UTF-8",
	}
	return code, msgInternal.SendEmail(email)
}
