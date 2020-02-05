package request

import (
	"time"

	"github.com/chadhao/logit/config"
	"github.com/chadhao/logit/modules/user/model"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	RefreshTokenRequest struct {
		Toekn string `json:"token"`
	}
	LoginRequest struct {
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Licence  string `json:"licence"`
		Password string `json:"password"`
	}
	UserRegistrationRequest struct {
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	DriverRegistrationRequest struct {
		Id            primitive.ObjectID `json:"id"`
		LicenceNumber string             `json:"licenceNumber"`
		DateOfBirth   time.Time          `json:"dateOfBirth"`
		Firstnames    string             `json:"firstnames"`
		Surname       string             `json:"surname"`
	}
	ExistanceRequest struct {
		Phone   string `json:"phone"`
		Email   string `json:"email"`
		Licence string `json:"licence"`
	}
)

func (r *RefreshTokenRequest) Validate(c config.Config) (*model.User, error) {
	u := model.User{}

	key, _ := c.Get("system.jwt.refresh.key")
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return key, nil
	}
	token, err := jwt.Parse(r.Toekn, keyFunc)
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.StandardClaims)
	userId, err := primitive.ObjectIDFromHex(claims.Subject)
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

func (r *UserRegistrationRequest) Reg() (*model.User, error) {
	// Should add Request content validation here
	u := model.User{
		Phone:    r.Phone,
		Email:    r.Email,
		Password: r.Password,
	}

	if err := u.Create(); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *DriverRegistrationRequest) Reg() (*model.Driver, error) {
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
