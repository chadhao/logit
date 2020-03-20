package request

import (
	"errors"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UserRegRequest struct {
		Phone    string `json:"phone" valid:"numeric,stringlength(8|11)"`
		Code     string `json:"code" valid:"numeric"`
		Email    string `json:"email" valid:"email"`
		Password string `json:"password" valid:"stringlength(6|32)"`
	}
	UserUpdateRequest struct {
		Password string `json:"password" valid:"stringlength(6|32)"`
	}
	DriverRegRequest struct {
		ID            primitive.ObjectID `json:"id"`
		LicenseNumber string             `json:"licenseNumber"`
		DateOfBirth   time.Time          `json:"dateOfBirth"`
		Firstnames    string             `json:"firstnames"`
		Surname       string             `json:"surname"`
	}
	TransportOperatorRegRequest struct {
		ID            primitive.ObjectID `json:"id"`
		LicenseNumber string             `json:"licenseNumber"`
		Name          string             `json:"name"`
	}
)

func (r *UserRegRequest) Reg() (*model.User, error) {
	if _, err := valid.ValidateStruct(r); err != nil {
		return nil, err
	}

	red := model.Redis{Key: r.Phone}
	if code, err := red.Get(); err != nil || r.Code != code {
		return nil, errors.New("verification code does not match")
	}

	u := model.User{
		Phone:     r.Phone,
		Email:     r.Email,
		Password:  r.Password,
		CreatedAt: time.Now(),
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
		ID:            r.ID,
		LicenseNumber: r.LicenseNumber,
		DateOfBirth:   r.DateOfBirth,
		Firstnames:    r.Firstnames,
		Surname:       r.Surname,
		CreatedAt:     time.Now(),
	}

	if err := d.Create(); err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *TransportOperatorRegRequest) Reg() (*model.TransportOperator, error) {
	// Should add Request content validation here
	d := model.TransportOperator{
		ID:            r.ID,
		LicenseNumber: r.LicenseNumber,
		Name:          r.Name,
		IsVerified:    false,
		SuperIDs:      []primitive.ObjectID{r.ID},
		CreatedAt:     time.Now(),
	}

	if err := d.Create(); err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *UserUpdateRequest) Replace(user *model.User) (err error) {
	if _, err := valid.ValidateStruct(r); err != nil {
		return err
	}
	user.Password = r.Password
	return nil
}
