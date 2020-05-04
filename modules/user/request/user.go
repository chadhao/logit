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
		Password *string `json:"password,omitempty"`
		Pin      *string `json:"pin,omitempty"`
	}
	DriverRegRequest struct {
		ID            primitive.ObjectID `json:"id"`
		LicenseNumber string             `json:"licenseNumber" valid:"stringlength(5|8)`
		DateOfBirth   time.Time          `json:"dateOfBirth" valid:"required"`
		Firstnames    string             `json:"firstnames" valid:"required"`
		Surname       string             `json:"surname" valid:"required"`
		Pin           string             `json:"pin" valid:"stringlength(4|4)`
	}
	TransportOperatorRegRequest struct {
		ID            primitive.ObjectID `json:"id"`
		LicenseNumber string             `json:"licenseNumber"`
		IsCompany     bool               `json:"isCompany"`
		Name          string             `json:"name"`
		Contact       *string            `json:"contact,omitempty"`
	}
	TransportOperatorUpdateRequest struct {
		ID            primitive.ObjectID `json:"id"`
		LicenseNumber string             `json:"licenseNumber"`
		Name          string             `json:"name"`
	}
	TransportOperatorAddIdentityRequest struct {
		TransportOperatorID primitive.ObjectID `json:"transportOperatorID"`
		UserID              primitive.ObjectID `json:"userID"`
		Identity            model.TOIdentity   `json:"identity"`
		Contact             *string            `json:"contact,omitempty"`
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
	if _, err := valid.ValidateStruct(r); err != nil {
		return nil, err
	}
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

func (r *TransportOperatorRegRequest) Reg(uid primitive.ObjectID) (*model.TransportOperator, error) {
	// Should add Request content validation here
	if r.IsCompany && r.Contact == nil {
		return nil, errors.New("contact is required")
	}

	d := model.TransportOperator{
		ID:            primitive.NewObjectID(),
		LicenseNumber: r.LicenseNumber,
		Name:          r.Name,
		IsVerified:    false,
		IsCompany:     r.IsCompany,
		CreatedAt:     time.Now(),
	}

	if err := d.Create(); err != nil {
		return nil, err
	}

	if _, err := d.AddIdentity(uid, model.TO_SUPER, r.Contact); err != nil {
		d.Delete()
		return nil, err
	}
	// 自雇性质注册时，当已经有driver信息时，自动添加为TO下driver
	if !d.IsCompany {
		driver := model.Driver{ID: uid}
		if err := driver.Find(); err == nil {
			d.AddIdentity(uid, model.TO_DRIVER, nil)
		}
	}

	return &d, nil
}

func (r *TransportOperatorUpdateRequest) Update() (*model.TransportOperator, error) {
	d := model.TransportOperator{
		ID: r.ID,
	}

	if err := d.Find(); err != nil {
		return nil, err
	}

	if r.LicenseNumber != "" {
		d.LicenseNumber = r.LicenseNumber
	}
	if r.Name != "" {
		d.Name = r.Name
	}

	if err := d.Update(); err != nil {
		return nil, err
	}
	return &d, nil
}

// Add .
func (r *TransportOperatorAddIdentityRequest) Add() (*model.TransportOperatorIdentity, error) {
	if r.Identity != model.TO_SUPER && r.Identity != model.TO_ADMIN {
		return nil, errors.New("identity not valid")
	}
	d := model.TransportOperator{
		ID: r.TransportOperatorID,
	}

	if err := d.Find(); err != nil {
		return nil, err
	}
	if !d.IsCompany || !d.IsVerified {
		return nil, errors.New("can only add identity to verified company")
	}

	identity, err := d.AddIdentity(r.UserID, r.Identity, r.Contact)
	if err != nil {
		return nil, err
	}
	return identity, nil
}

func (r *UserUpdateRequest) Replace(user *model.User) (err error) {
	if (UserUpdateRequest{}) == *r {
		return errors.New("at least one update request is required")
	}
	if r.Password != nil {
		if len(*r.Password) < 6 || len(*r.Password) > 32 {
			return errors.New("the length of password should be between 6 to 32")
		}
		if user.Password == *r.Password {
			return errors.New("password is the same as the old one")
		}
		user.Password = *r.Password
	}
	if r.Pin != nil {
		if len(*r.Pin) != 4 {
			return errors.New("the pin should be 4 digits")
		}
		if user.Pin == *r.Pin {
			return errors.New("pin is the same as the old one")
		}
		user.Pin = *r.Pin
	}
	return nil
}
