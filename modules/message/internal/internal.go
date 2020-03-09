package internal

import "github.com/chadhao/logit/modules/message/model"

// SendTxt 发送信息
func SendTxt(txt model.Txt) error {
	return txt.Send()
}

// SendEmail 发送信息
func SendEmail(email model.Email) error {
	return email.Send()
}
