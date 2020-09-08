package internals

import (
	"errors"

	valid "github.com/asaskevich/govalidator"
	"github.com/chadhao/logit/modules/message/model"
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

// SendTxt 发送信息
func SendTxt(txt TxtRequest) error {
	if _, err := valid.ValidateStruct(txt); err != nil {
		return err
	}
	t := model.Txt{
		Number:  txt.Number,
		Message: txt.Message,
	}
	return t.Send()
}

// SendEmail 发送信息
func SendEmail(email EmailRequest) error {

	if _, err := valid.ValidateStruct(email); err != nil {
		return err
	}
	for _, v := range email.Recipients {
		if !valid.IsEmail(v) {
			return errors.New("not email")
		}
	}

	e := model.Email{
		Sender:     email.Sender,
		Recipients: email.Recipients,
		Subject:    email.Subject,
		HTMLBody:   email.HTMLBody,
		TextBody:   email.TextBody,
		CharSet:    email.CharSet,
	}
	return e.Send()
}
