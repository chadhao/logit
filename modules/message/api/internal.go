package api

import "github.com/chadhao/logit/modules/message/model"

// SendTxt 发送信息
func SendTxt(txt TxtRequest) error {
	t := model.Txt{
		Number:  txt.Number,
		Message: txt.Message,
	}
	return t.Send()
}

// SendEmail 发送信息
func SendEmail(email EmailRequest) error {
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
