package model

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
)

// Email 邮件
type Email struct {
	Sender     string   `json:"sender"`
	Recipients []string `json:"recipients"`
	Subject    string   `json:"subject"`
	HTMLBody   string   `json:"htmlBody"`
	TextBody   string   `json:"textBody"`
	CharSet    string   `json:"charSet"`
}

func (e *Email) valid() error {
	if _, err := awsSES.VerifyEmailAddress(&ses.VerifyEmailAddressInput{EmailAddress: aws.String(e.Sender)}); err != nil {
		return err
	}
	if e.CharSet == "" {
		e.CharSet = "UTF-8"
	}
	return nil
}

// Send 发送邮件
func (e *Email) Send() error {
	if err := e.valid(); err != nil {
		return err
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(e.Recipients),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(e.CharSet),
					Data:    aws.String(e.HTMLBody),
				},
				Text: &ses.Content{
					Charset: aws.String(e.CharSet),
					Data:    aws.String(e.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(e.CharSet),
				Data:    aws.String(e.Subject),
			},
		},
		Source: aws.String(e.Sender),
	}
	_, err := awsSES.SendEmail(input)
	return err
}
