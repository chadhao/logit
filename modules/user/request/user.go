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

func (r *TransportOperatorRegRequest) Reg(uid primitive.ObjectID) (*model.TransportOperator, error) {
	// Should add Request content validation here
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

	_, err := d.AddIdentity(uid, model.TO_SUPER, r.Contact)
	if err != nil {
		return nil, err
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

func (r *TransportOperatorAddIdentityRequest) Add() (*model.TransportOperatorIdentity, error) {
	d := model.TransportOperator{
		ID: r.TransportOperatorID,
	}

	if err := d.Find(); err != nil {
		return nil, err
	}

	identity, err := d.AddIdentity(r.UserID, r.Identity, r.Contact)
	if err != nil {
		return nil, err
	}
	return identity, nil
}

func (r *UserUpdateRequest) Replace(user *model.User) (err error) {
	if _, err := valid.ValidateStruct(r); err != nil {
		return err
	}
	user.Password = r.Password
	return nil
}
