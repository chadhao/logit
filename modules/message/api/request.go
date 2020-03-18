package api

import (
	"errors"

	valid "github.com/asaskevich/govalidator"
)

type (
	// TxtRequest .
	TxtRequest struct {
		Number  string `json:"number" valid:"numeric,stringlength(8|11)"`
		Message string `json:"message" valid:"stringlength(0|255)"`
	}
	// EmailRequest .
	EmailRequest struct {
		Sender     string   `json:"sender" valid:"email"`
		Recipients []string `json:"recipients" valid:"required"`
		Subject    string   `json:"subject" valid:"stringlength(0|50)"`
		HTMLBody   string   `json:"htmlBody" valid:"required"`
		TextBody   string   `json:"textBody" valid:"required"`
		CharSet    string   `json:"charSet" valid:"required"`
	}
)

func (t *TxtRequest) valid() error {
	_, err := valid.ValidateStruct(t)
	return err
}

func (e *EmailRequest) valid() error {
	_, err := valid.ValidateStruct(e)
	for _, v := range e.Recipients {
		if !valid.IsEmail(v) {
			return errors.New("not email")
		}
	}
	return err
}
